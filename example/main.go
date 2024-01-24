package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	freee "github.com/kurusugawa-computer/freee-go"
	"github.com/kurusugawa-computer/freee-go/oauth"
)

func main() {
	clientID := "xxxx"
	clientSecret := "xxxx"

	// アクセストークンを取得します。
	accessToken, err := freee.Authorize(clientID, clientSecret, 8080,
		oauth.WithPrompt(func(authorizeURL string) error {
			fmt.Println("次のURLにアクセスして認証を行ってください。")
			fmt.Println(authorizeURL)
			return nil
		}),
		oauth.WithRenderer(func(w http.ResponseWriter, authorizationCode string, err error) {
			content := "認証に成功しました。ブラウザを閉じてください。"
			if err != nil {
				content = "認証に失敗しました。ブラウザを閉じて、アプリケーションをもう一度実行してください。"
			}
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			w.Header().Set("Content-Length", strconv.Itoa(len(content)))
			w.WriteHeader(http.StatusOK)
			io.Copy(w, strings.NewReader(content))
		}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// freee API クライアントを作成します。
	client, err := freee.New(clientID, clientSecret, accessToken)
	if err != nil {
		log.Fatalln(err)
	}

	// ログインユーザーを取得します。
	loginUser, err := client.GetLoginUser()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(loginUser.ID)
}
