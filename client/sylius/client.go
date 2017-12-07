package sylius

import (
	"context"
	"fmt"
	"highlite-parser/client/sylius/transfer"
	"strings"
	"time"

	"net/http"

	"github.com/go-resty/resty"
)

const (
	tokenRefreshInterval     = 30 * time.Minute
	tokenRequestRetryTimeout = time.Second
	tokenRequestRetryCount   = 5

	requestRetryCount   int           = 3
	requestRetryTimeout time.Duration = 200 * time.Millisecond
	requestTimeout      time.Duration = time.Second
)

func NewClient(endpoint, clientID, clientSecret, username, password string) *client {
	client := &client{
		endpoint:     endpoint,
		clientID:     clientID,
		clientSecret: clientSecret,
		username:     username,
		password:     password,
	}

	go client.tokenServer()

	return client
}

type client struct {
	endpoint     string
	clientID     string
	clientSecret string
	username     string
	password     string

	tokenChan <-chan transfer.Token
}

// Gets tokenChan by username and password
func (s *client) getTokenByPassword(ctx context.Context) (*transfer.Token, error) {
	result := &transfer.Token{}
	resp, err := resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     s.clientID,
			"client_secret": s.clientSecret,
			"grant_type":    "password",
			"username":      s.username,
			"password":      s.password,
		}).
		SetResult(result).
		Post(s.getUrl("/oauth/v2/tokens"))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("getTokenByPassword: request is not OK, status is: %d", resp.StatusCode())
	}

	return result, nil
}

// Gets tokenChan by refresh tokenChan
func (s *client) getTokenByRefreshToken(ctx context.Context, refreshToken string) (*transfer.Token, error) {
	result := &transfer.Token{}
	resp, err := resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     s.clientID,
			"client_secret": s.clientSecret,
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
		}).
		SetResult(result).
		Post(s.getUrl("/oauth/v2/tokens"))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("getTokenByRefreshToken: request is not OK, status is: %d", resp.StatusCode())
	}

	return result, nil
}

// Token delivery server. Gets tokens from Sylius OAuth server and writes them to a channel.
func (s *client) tokenServer() {
	tokenChan := make(chan transfer.Token)
	defer close(tokenChan)

	s.tokenChan = tokenChan

	var token *transfer.Token
	var refreshToken <-chan time.Time

	obtainToken := make(chan bool, 1)
	obtainToken <- true

	for keepRunning := true ; keepRunning ; {
		var tokenRequestChan chan transfer.Token
		if token != nil {
			tokenRequestChan = tokenChan
		}

		select {
		case tokenRequestChan <- *token:
			fmt.Println("wrote token to token channel")

		case <-obtainToken:
			newToken, err := s.obtainTokenByPasswordAndUsername()
			if err != nil {
				keepRunning = false
			} else {
				token = newToken
				refreshToken = time.After(tokenRefreshInterval)
			}

		case <-refreshToken:
			newToken, err := s.getTokenByRefreshToken(s.getContextWithTimeout(), token.RefreshToken)
			if err != nil {
				obtainToken <- true
			} else {
				token = newToken
				refreshToken = time.After(tokenRefreshInterval)
			}

		}
	}

	close(tokenChan)
}

// Tries to get tokenChan by username and password.
// Retries for a const amount of times if api request fails.
func (s *client) obtainTokenByPasswordAndUsername() (*transfer.Token, error) {
	for i := 0; i < tokenRequestRetryCount; i++ {
		token, err := s.getTokenByPassword(s.getContextWithTimeout())
		if err != nil {
			time.Sleep(tokenRequestRetryTimeout)
		} else {
			return token, nil
		}
	}

	return nil, fmt.Errorf("can't obtain tokenChan by password and username after %d retries", tokenRequestRetryTimeout)
}

// Gets a tokenChan structure.
func (s *client) getToken(updateExisting bool) (*transfer.Token, error) {
	select {
	case token, ok := <-s.tokenChan:
		if !ok {
			return nil, fmt.Errorf("can't get tokenChan: ")
		} else {
			return &token, nil
		}

	case <-time.After(4 * time.Second):
		return nil, fmt.Errorf("get tokenChan timeout")
	}
}

// Returns full url using endpoint and resource paths.
func (s *client) getUrl(path string) string {
	return strings.TrimSuffix(s.endpoint, "/") + "/" + strings.TrimPrefix(path, "/")
}

// Gets context with timeout
func (s *client) getContextWithTimeout() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	return ctx
}
