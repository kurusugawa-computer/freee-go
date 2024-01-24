package freee

import (
	"net/http"
	"sync"
	"time"

	"github.com/kurusugawa-computer/freee-go/token"
)

type AccessToken = token.TokenInfo

func newTokenManager(clientID string, clientSecret string, accessToken *AccessToken, httpClient *http.Client) *tokenManager {
	return &tokenManager{
		clientID:       clientID,
		clientSecret:   clientSecret,
		token:          accessToken,
		httpClient:     http.DefaultClient,
		mutex:          sync.Mutex{},
		onRefreshToken: nil,
	}
}

type tokenManager struct {
	clientID       string
	clientSecret   string
	token          *AccessToken
	httpClient     *http.Client
	mutex          sync.Mutex
	onRefreshToken func(*AccessToken) error
}

func (m *tokenManager) OnRefreshToken(f func(*AccessToken) error) {
	m.onRefreshToken = f
}

func (m *tokenManager) GetAccessToken() (*AccessToken, error) {
	m.mutex.Lock()
	accessToken := m.token

	if time.Now().After(time.Unix(accessToken.CreatedAt, 0).Add(time.Duration(accessToken.ExpiresIn) * time.Second)) {
		accessToken, err := token.RefreshAccessToken(m.clientID, m.clientSecret, accessToken.RefreshToken, token.WithHTTPClient(m.httpClient))
		if err != nil {
			return nil, err
		}
		if m.onRefreshToken != nil {
			if err := m.onRefreshToken(accessToken); err != nil {
				return nil, err
			}
		}
		m.token = accessToken
	}

	m.mutex.Unlock()
	return accessToken, nil
}
