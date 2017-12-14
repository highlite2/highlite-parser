package sylius

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty"

	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius/transfer"
)

const (
	tokenRefreshInterval   = 30 * time.Minute
	tokenRequestRetryCount = 10

	requestTimeout time.Duration = time.Second

	methodGet  string = "get"
	methodPost string = "post"
)

// IClient is a Sylius client interface
type IClient interface {
	GetTaxon(ctx context.Context, code string) (*transfer.Taxon, error)
}

// NewClient is a Sylius client constructor.
func NewClient(logger log.ILogger, endpoint string, auth Auth) IClient {
	c := &client{
		endpoint:  endpoint,
		auth:      auth,
		logger:    logger,
		tokenChan: make(chan *transfer.Token),
	}

	go c.tokenServer()

	c.getToken()

	return c
}

// Auth contains credentials to obtain Sylius API token.
type Auth struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

type client struct {
	endpoint  string
	auth      Auth
	logger    log.ILogger
	tokenChan chan *transfer.Token
}

// Gets tokenChan by Username and Password.
func (c *client) getTokenByPassword(ctx context.Context) (*transfer.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	url := c.getURL("/oauth/v2/token")

	result := &transfer.Token{}
	resp, err := resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     c.auth.ClientID,
			"client_secret": c.auth.ClientSecret,
			"grant_type":    "password",
			"username":      c.auth.Username,
			"password":      c.auth.Password,
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

// Gets tokenChan by refresh tokenChan.
func (c *client) getTokenByRefreshToken(ctx context.Context, refreshToken string) (*transfer.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	url := c.getURL("/oauth/v2/token")

	result := &transfer.Token{}
	resp, err := resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     c.auth.ClientID,
			"client_secret": c.auth.ClientSecret,
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
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

// Token delivery server. Gets tokens from Sylius OAuth server and writes them to a channel.
// Makes a token background update.
func (c *client) tokenServer() {
	var token *transfer.Token

	refreshToken := make(<-chan time.Time)
	obtainToken := make(chan bool, 1)
	obtainToken <- true

	c.logger.Debug("Starting token delivery server")
	for keepRunning := true; keepRunning; {
		var tokenRequestChan chan *transfer.Token
		if token != nil {
			tokenRequestChan = c.tokenChan
		}

		select {
		case tokenRequestChan <- token:

		case <-obtainToken:
			c.logger.Debug("Trying to obtain token by password")
			newToken, err := c.obtainTokenByPasswordAndUsername()
			if err != nil {
				c.logger.Errorf("Can't get token: %s", err.Error())
				keepRunning = false
			} else {
				c.logger.Debug("Successfully received token")
				token = newToken
				refreshToken = time.After(tokenRefreshInterval)
			}

		case <-refreshToken:
			newToken, err := c.getTokenByRefreshToken(context.Background(), token.RefreshToken)
			if err != nil {
				c.logger.Errorf("Can't refresh token using refresh token: %s", err.Error())
				obtainToken <- true
				token = nil
			} else {
				c.logger.Debug("Successfully updated token")
				token = newToken
				refreshToken = time.After(tokenRefreshInterval)
			}

		}
	}

	c.logger.Debug("Stopping token delivery server")
}

// Tries to get tokenChan by Username and Password.
// Retries for a const amount of times if api request fails.
func (c *client) obtainTokenByPasswordAndUsername() (*transfer.Token, error) {
	for i := 0; i < tokenRequestRetryCount; i++ {
		token, err := c.getTokenByPassword(context.Background())
		if err != nil {
			c.logger.Warnf("Failed to obtain password for the %d time: %s", i+1, err.Error())
			time.Sleep(time.Second)
		} else {
			return token, nil
		}
	}

	return nil, fmt.Errorf("can't obtain token by password and username after %d retries", tokenRequestRetryCount)
}

// Gets a token.
func (c *client) getToken() (string, error) {
	token, ok := <-c.tokenChan
	if !ok {
		return "", fmt.Errorf("can't get token: token chan was closed")
	}

	return token.AccessToken, nil
}

// Returns full url using endpoint and resource paths.
func (c *client) getURL(path string, args ...interface{}) string {
	if len(args) > 0 {
		path = fmt.Sprintf(path, args...)
	}

	return strings.TrimSuffix(c.endpoint, "/") + "/" + strings.TrimPrefix(path, "/")
}

// Performs GET request
func (c *client) requestGet(ctx context.Context, url string, result interface{}) error {
	return c.request(ctx, methodGet, url, result)
}

// Performs a request. Sets authorization token and handles errors.
// Creates context with timeout.
func (c *client) request(ctx context.Context, method string, url string, result interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	var err error

	token, err := c.getToken()
	if err != nil {
		c.logger.Errorf("Failed to get token during [%s] %s request", method, url)

		return err
	}

	request := resty.R().SetContext(ctx).SetHeader("Authorization", "Bearer "+token).SetResult(result)

	var res *resty.Response

	switch method {
	case methodGet:
		res, err = request.Get(url)
	case methodPost:
		res, err = request.Post(url)
	default:
		err = fmt.Errorf("unknown method")
	}

	if err != nil {
		c.logger.Errorf("Request to [%s] %s failed: %s", method, url, err.Error())

		return fmt.Errorf("request to [%s] %s failed: %s", method, url, err.Error())
	}

	if res.StatusCode() != http.StatusOK {
		c.logger.Errorf("Request to [%s] %s ended with status %s", method, url, res.Status())

		return fmt.Errorf("request to [%s] %s ended with status %s", method, url, res.Status())
	}

	return nil
}
