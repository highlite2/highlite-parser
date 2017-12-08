package sylius

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"highlite-parser/client/sylius/transfer"
	"highlite-parser/log"

	"github.com/go-resty/resty"
)

const (
	tokenRefreshInterval     = 30 * time.Minute
	tokenRequestRetryTimeout = time.Second
	tokenRequestRetryCount   = 5

	requestTimeout time.Duration = time.Second
)

func NewClient(log log.Logger, endpoint string, auth Auth) *client {
	client := &client{
		endpoint: endpoint,
		auth:     auth,
		log:      log,
	}

	go client.tokenServer()

	return client
}

type Auth struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

type client struct {
	endpoint  string
	auth      Auth
	log       log.Logger
	tokenChan <-chan *transfer.Token
}

// Gets tokenChan by Username and Password
func (s *client) getTokenByPassword(ctx context.Context) (*transfer.Token, error) {
	url := s.getUrl("/oauth/v2/token")
	result := &transfer.Token{}
	resp, err := resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     s.auth.ClientID,
			"client_secret": s.auth.ClientSecret,
			"grant_type":    "password",
			"username":      s.auth.Username,
			"password":      s.auth.Password,
		}).
		SetResult(result).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %s", url, err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request to %s ended with status %s", url, resp.Status())
	}

	return result, nil
}

// Gets tokenChan by refresh tokenChan
func (s *client) getTokenByRefreshToken(ctx context.Context, refreshToken string) (*transfer.Token, error) {
	url := s.getUrl("/oauth/v2/token")
	result := &transfer.Token{}
	resp, err := resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     s.auth.ClientID,
			"client_secret": s.auth.ClientSecret,
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
		}).
		SetResult(result).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %s", url, err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request to %s ended with status %d", url, resp.StatusCode())
	}

	return result, nil
}

// Token delivery server. Gets tokens from Sylius OAuth server and writes them to a channel.
// Makes a background update of tokens.
func (s *client) tokenServer() {
	tokenChan := make(chan *transfer.Token)
	defer close(tokenChan)

	s.tokenChan = tokenChan

	var token *transfer.Token
	var refreshToken <-chan time.Time

	obtainToken := make(chan bool, 1)
	obtainToken <- true

	s.log.Info("Starting token delivery server")

	for keepRunning := true; keepRunning; {
		var tokenRequestChan chan *transfer.Token
		if token != nil {
			tokenRequestChan = tokenChan
		}

		select {
		case tokenRequestChan <- token:

		case <-obtainToken:
			s.log.Info("Trying to obtain token by password")
			newToken, err := s.obtainTokenByPasswordAndUsername()
			if err != nil {
				s.log.Errorf("Can't get token: %s", err.Error())
				keepRunning = false
			} else {
				s.log.Info("Successfully received token")
				token = newToken
				refreshToken = time.After(tokenRefreshInterval)
			}

		case <-refreshToken:
			newToken, err := s.getTokenByRefreshToken(s.getContextWithTimeout(), token.RefreshToken)
			if err != nil {
				s.log.Errorf("Can't refresh token using refresh token: %s", err.Error())
				obtainToken <- true
				token = nil
			} else {
				s.log.Info("Successfully updated token")
				token = newToken
				refreshToken = time.After(tokenRefreshInterval)
			}

		}
	}

	s.log.Info("Stopping token delivery server")
}

// Tries to get tokenChan by Username and Password.
// Retries for a const amount of times if api request fails.
func (s *client) obtainTokenByPasswordAndUsername() (*transfer.Token, error) {
	for i := 0; i < tokenRequestRetryCount; i++ {
		token, err := s.getTokenByPassword(s.getContextWithTimeout())
		if err != nil {
			s.log.Warnf("Failed to obtain password for the %d time: %s", i+1, err.Error())
			time.Sleep(tokenRequestRetryTimeout)
		} else {
			return token, nil
		}
	}

	return nil, fmt.Errorf("can't obtain token by password and username after %d retries", tokenRequestRetryCount)
}

// Gets a tokenChan structure.
func (s *client) getToken(updateExisting bool) (string, error) {
	select {
	case token, ok := <-s.tokenChan:
		if !ok {
			return "", fmt.Errorf("can't get token: token chan is closed")
		} else {
			return token.AccessToken, nil
		}

	case <-time.After(time.Second):
		return "", fmt.Errorf("can't get token: timeout")
	}
}

// Returns full url using Endpoint and resource paths.
func (s *client) getUrl(path string) string {
	return strings.TrimSuffix(s.endpoint, "/") + "/" + strings.TrimPrefix(path, "/")
}

// Gets context with timeout
func (s *client) getContextWithTimeout() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	return ctx
}
