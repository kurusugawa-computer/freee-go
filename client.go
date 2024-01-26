package freee

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

type opt struct {
	HTTPClient    *http.Client
	beforeRequest []func(*http.Request) (*http.Request, error)
	afterResponse []func(*http.Response) (*http.Response, error)
}

type OptFunc func(*opt)

func WithHTTPClient(httpClient *http.Client) func(*opt) {
	return func(o *opt) {
		o.HTTPClient = httpClient
	}
}

type Hooks struct {
	BeforeRequest func(*http.Request) (*http.Request, error)
	AfterResponse func(*http.Response) (*http.Response, error)
}

func WithHooks(hooks Hooks) func(*opt) {
	return func(o *opt) {
		if hooks.BeforeRequest != nil {
			o.beforeRequest = append(o.beforeRequest, hooks.BeforeRequest)
		}
		if hooks.AfterResponse != nil {
			o.afterResponse = append(o.afterResponse, hooks.AfterResponse)
		}
	}
}

func New(clientID string, clientSecret string, accessToken *AccessToken, opts ...OptFunc) (*Client, error) {
	if accessToken == nil {
		return nil, errors.New("access token is nil")
	}

	o := &opt{
		HTTPClient: http.DefaultClient,
	}
	for _, of := range opts {
		of(o)
	}

	c := &Client{
		httpClient:    o.HTTPClient,
		Token:         newTokenManager(clientID, clientSecret, accessToken, o.HTTPClient),
		beforeRequest: o.beforeRequest,
		afterResponse: o.afterResponse,
	}

	return c, nil
}

type Client struct {
	httpClient    *http.Client
	Token         *tokenManager
	beforeRequest []func(*http.Request) (*http.Request, error)
	afterResponse []func(*http.Response) (*http.Response, error)
}

type response struct {
	*http.Response
}

func (r *response) Close() error {
	_, err := io.Copy(io.Discard, r.Body)
	if err1 := r.Body.Close(); err == nil {
		err = err1
	}
	return err
}

func (r *response) Parse(v any) error {
	defer r.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid body: %v", err)
	}
	return nil
}

func (c *Client) do(method string, targetURL string, query url.Values, payload any) (*response, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid request url: %v", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var body io.Reader
	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("invalid request payload: %v", err)
		}
		body = bytes.NewReader(buf)
	}

	accessToken, err := c.Token.GetAccessToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}

	for _, hook := range c.beforeRequest {
		req, err = hook(req)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	for _, hook := range c.afterResponse {
		resp, err = hook(resp)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return &response{resp}, nil
	}

	return nil, handleError(&response{resp})
}

func handleError(r *response) error {
	defer r.Close()

	invalidStatusCode := func() error {
		return errors.New("invalid status code: " + r.Status)
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return invalidStatusCode()
	}

	switch mediaType {
	case "application/problem+json":
		problem := struct {
			StatusCode int `json:"status_code"`
			Errors     []struct {
				Type     string   `json:"type"`
				Messages []string `json:"messages"`
			} `json:"errors"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
			return invalidStatusCode()
		}
		errs := make([]error, 0, len(problem.Errors))
		for _, err := range problem.Errors {
			errs = append(errs, errors.New(strings.Join(err.Messages, " ")+" ("+err.Type+")"))
		}
		return errors.Join(errs...)

	case "application/json":
		errResponse := struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&errResponse); err != nil {
			return invalidStatusCode()
		}
		return errors.New(errResponse.ErrorDescription + " (" + errResponse.Error + ")")
	}

	return invalidStatusCode()
}
