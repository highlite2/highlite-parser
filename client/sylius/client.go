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
	tokenRequestRetryCount = 5
	tokenRequestRetryTimeout = 200 * time.Millisecond

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
		tokenRequest: make(chan *tokenRequest),
	}

	go client.tokenDeliveryServer()

	return client
}

type client struct {
	endpoint     string
	clientID     string
	clientSecret string
	username     string
	password     string

	tokenRequest chan *tokenRequest
}

type tokenRequest struct {
	update bool
	token  chan transfer.Token
	error  chan error
}

func (s *client) GetCategories() {

}

// Gets token by username and password
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
		return nil, fmt.Errorf("request is not ok, status is: %d", resp.StatusCode())
	}

	return result, nil
}

//func (s *client) tokenUpdateServer() {
//	// first we need to obtain token by username and password
//	var token *transfer.Token
//	for ; ; token == nil {
//
//	}
//
//}

// Token Delivery server. Gets tokens from Sylius OAuth server and responses for token requests.
func (s *client) tokenDeliveryServer() {
	var token *transfer.Token

	for request := range s.tokenRequest {
		var err error

		if request.update {
			token = nil
		}

		if token == nil {
			for i := requestRetryCount; token == nil && i > 0; i-- {
				ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
				token, err = s.getTokenByPassword(ctx)
				if err != nil && i > 0 {
					time.Sleep(requestRetryTimeout)
				}
			}
		}

		if err != nil {
			request.error <- err
		} else {
			request.token <- *token
		}
	}
}

// Gets a token structure.
func (s *client) getToken(updateExisting bool) (*transfer.Token, error) {
	request := &tokenRequest{
		update: updateExisting,
		token:  make(chan transfer.Token),
		error:  make(chan error),
	}

	s.tokenRequest <- request

	select {
	case token := <-request.token:
		return &token, nil
	case err := <-request.error:
		return nil, fmt.Errorf("get token error: %s", err.Error())
	case <-time.After(4 * time.Second):
		return nil, fmt.Errorf("get token timeout")
	}
}

// Returns full url using endpoint and resource paths.
func (s *client) getUrl(path string) string {
	return strings.TrimSuffix(s.endpoint, "/") + "/" + strings.TrimPrefix(path, "/")
}
