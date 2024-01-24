package token

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
)

type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // トークンが有効な CreatedAt からの秒数
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"` // トークンが作成された Unix 秒
	CompanyID    int    `json:"company_id"`
}

type opt struct {
	HTTPClient *http.Client
}

type OptFunc func(*opt)

func WithHTTPClient(client *http.Client) func(*opt) {
	return func(o *opt) {
		o.HTTPClient = client
	}
}

// GetAccessToken は認可コードを使用して freee の API アクセスに必要なトークン情報を取得します。
func GetAccessToken(clientID string, clientSecret string, redirectURI string, authorizeCode string, opts ...OptFunc) (*TokenInfo, error) {
	q := url.Values{}
	q.Set("grant_type", "authorization_code")
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("code", authorizeCode)
	q.Set("redirect_uri", redirectURI)
	return requestToken(q, opts...)
}

// RefreshAccessToken はリフレッシュトークンを使用して freee の API アクセスに必要なトークン情報を取得します。
func RefreshAccessToken(clientID string, clientSecret string, refreshToken string, opts ...OptFunc) (*TokenInfo, error) {
	q := url.Values{}
	q.Set("grant_type", "refresh_token")
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("refresh_token", refreshToken)
	return requestToken(q, opts...)
}

func requestToken(q url.Values, opts ...OptFunc) (*TokenInfo, error) {
	o := &opt{
		HTTPClient: http.DefaultClient,
	}
	for _, of := range opts {
		of(o)
	}

	u, _ := url.Parse("https://accounts.secure.freee.co.jp/public_api/token")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	switch {
	case resp.StatusCode == 200:
		token := &TokenInfo{}
		if tError := json.NewDecoder(resp.Body).Decode(token); tError != nil {
			return nil, errors.New("invalid body: " + tError.Error())
		}
		return token, nil

	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		if mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type")); mediaType != "application/json" || err != nil {
			return nil, errors.New("invalid status code: " + resp.Status)
		}
		errResponse := &struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(errResponse); err != nil {
			return nil, errors.New("invalid status code: " + resp.Status)
		}
		return nil, errors.New(errResponse.ErrorDescription + " (" + errResponse.Error + ")")

	default:
		return nil, errors.New("invalid status code: " + resp.Status)
	}
}
