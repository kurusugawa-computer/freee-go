package oauth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type opt struct {
	Prompt   func(string) error
	Renderer func(http.ResponseWriter, string, error)
}

type OptFunc func(*opt)

func WithPrompt(prompt func(string) error) func(*opt) {
	return func(o *opt) {
		o.Prompt = prompt
	}
}

func WithRenderer(renderer func(http.ResponseWriter, string, error)) func(*opt) {
	return func(o *opt) {
		o.Renderer = renderer
	}
}

// Authorize は freee の OAuth 認証を実行し、認可コードを取得します。
//
// この関数は次の手順で OAuth 認証を行います。
//  1. freee の OAUTH 認証のために一時的な HTTP サーバーを起動します。
//  2. ユーザーに freee の認証 URL へアクセスすることを促します。
//     ※ WithPrompt でこの処理を変更できます。
//  3. ユーザーが freee の認証 URL へアクセスし、ユーザー認証を行い、アプリの利用を認可すると
//     一時的に起動した HTTP サーバーへリダイレクトされます。
//     ※ リダイレクト先のレンダリング内容は WithRenderer で変更できます。
//  4. 一時的に起動した HTTP サーバーはリダイレクトされた情報を検証し、認可コードを取得します。
//
// この関数で認証を行うには、認証したいアプリの「コールバックURL」が
// http://localhost:<port>/ と完全に一致している必要があります。
func Authorize(clientID string, port int, opts ...OptFunc) (string, error) {
	o := &opt{
		Prompt: func(aURL string) error {
			fmt.Println("次のURLにアクセスして認証してください。")
			fmt.Println(aURL)
			return nil
		},
		Renderer: func(aWriter http.ResponseWriter, aAuthorizationCode string, aError error) {
			tContent := "認証に成功しました。ブラウザを閉じてください。"
			if aError != nil {
				tContent = "認証に失敗しました。ブラウザを閉じて、アプリケーションをもう一度実行してください。"
			}
			aWriter.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			aWriter.Header().Set("Content-Length", strconv.Itoa(len(tContent)))
			aWriter.WriteHeader(http.StatusOK)
			io.Copy(aWriter, strings.NewReader(tContent))
		},
	}
	for _, of := range opts {
		of(o)
	}

	state := generateRandomString(32)
	authorizeURL := makeAuthorizeURL(clientID, port, state)

	if err := o.Prompt(authorizeURL); err != nil {
		return "", nil
	}

	type Result struct {
		AuthorizationCode string
		Error             error
	}
	resultCh := make(chan Result)
	defer close(resultCh)

	go func() {
		shutdownCh := make(chan struct{})
		defer close(shutdownCh)

		authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()

			tAuthorizationCode := q.Get("code")

			var err error
			if q.Get("state") != state {
				tAuthorizationCode = ""
				err = errors.New("state mismatch")
			}

			if q.Get("error") != "" {
				err = errors.New(q.Get("error_description"))
			}

			o.Renderer(w, tAuthorizationCode, err)

			resultCh <- Result{tAuthorizationCode, err}
			shutdownCh <- struct{}{}
		})

		notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not Found", http.StatusNotFound)
		})

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				authHandler.ServeHTTP(w, r)
			default:
				notFoundHandler.ServeHTTP(w, r)
			}
		})

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		}
		go func() {
			<-shutdownCh
			server.Shutdown(context.Background())
		}()
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				// サーバー停止以外のエラーをエラーとする
				resultCh <- Result{"", err}
			}
		}
	}()

	result := <-resultCh

	return result.AuthorizationCode, result.Error
}

func makeAuthorizeURL(clientID string, port int, state string) string {
	u, _ := url.Parse("https://accounts.secure.freee.co.jp/public_api/authorize")
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", fmt.Sprintf("http://localhost:%d/", port)) // アプリの「コールバックURL」と文字列的に一致する必要がある
	q.Set("state", state)
	q.Set("prompt", "select_company")
	u.RawQuery = q.Encode()
	return u.String()
}

func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = letters[rand.Intn(len(letters))]
	}
	return string(buf)
}
