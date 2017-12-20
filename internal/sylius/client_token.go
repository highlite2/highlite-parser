package sylius

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty"

	"highlite-parser/internal/sylius/transfer"
)

// Gets a token.
func (c *Client) getToken() (string, error) {
	token, ok := <-c.tokenChan
	if !ok {
		return "", fmt.Errorf("can't get token: token chan was closed")
	}

	return token.AccessToken, nil
}

// Gets tokenChan by Username and Password.
func (c *Client) getTokenByPassword(ctx context.Context) (*transfer.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
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
func (c *Client) getTokenByRefreshToken(ctx context.Context, refreshToken string) (*transfer.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
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
func (c *Client) tokenServer() {
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

	close(c.tokenChan)

	c.logger.Debug("Stopping token delivery server")
}

// Tries to get tokenChan by Username and Password.
// Retries for a const amount of times if api request fails.
func (c *Client) obtainTokenByPasswordAndUsername() (*transfer.Token, error) {
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
