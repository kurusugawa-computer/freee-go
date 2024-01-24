package freee

import (
	"fmt"

	"github.com/kurusugawa-computer/freee-go/oauth"
	"github.com/kurusugawa-computer/freee-go/token"
)

// Authorize は freee の OAuth 認証を実行し、アクセストークンを取得します。
//
// この関数は次の手順で OAuth 認証を行います。
//  1. freee の OAUTH 認証のために一時的な HTTP サーバーを起動します。
//  2. ユーザーに freee の認証 URL へアクセスすることを促します。
//     ※ WithPrompt でこの処理を変更できます。
//  3. ユーザーが freee の認証 URL へアクセスし、ユーザー認証を行い、アプリの利用を認可すると
//     一時的に起動した HTTP サーバーへリダイレクトされます。
//     ※ リダイレクト先のレンダリング内容は WithRenderer で変更できます。
//  4. 一時的に起動した HTTP サーバーはリダイレクトされた情報を検証し、認可コードを取得します。
//  5. 認可コードを使用して、アクセストークンを取得します。
//
// この関数で認証を行うには、認証したいアプリの「コールバックURL」が
// http://localhost:<port>/ と完全に一致している必要があります。
func Authorize(clientID string, clientSecret string, callbackPort int, opts ...oauth.OptFunc) (*AccessToken, error) {
	authorizationCode, err := oauth.Authorize(clientID, callbackPort, opts...)
	if err != nil {
		return nil, err
	}

	redirectURI := fmt.Sprintf("http://localhost:%d/", callbackPort)
	accessToken, err := token.GetAccessToken(clientID, clientSecret, redirectURI, authorizationCode)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
