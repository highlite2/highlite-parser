package highlite

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"highlite2-import/internal/log"

	"github.com/go-resty/resty"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// NewClient gets new highlite client.
func NewClient(logger log.ILogger, login, password, loginEndpoint, itemsEndpoint string) *Client {
	client := resty.New()
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(5))

	return &Client{
		logger: logger,

		login:         login,
		password:      password,
		loginEndpoint: loginEndpoint,
		itemsEndpoint: itemsEndpoint,

		client: client,
	}
}

// Client is a highlite client.
type Client struct {
	logger log.ILogger

	login         string
	password      string
	loginEndpoint string
	itemsEndpoint string

	client *resty.Client
}

// GetItemsReader returns a Reader instance with highlite items.
func (r *Client) GetItemsReader(ctx context.Context) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	cookies, err := r.getCookies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session cookies: %s", err.Error())
	}

	r.client.SetCookies(cookies)
	resp, err := r.client.R().SetContext(ctx).Get(r.itemsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get highlite items: %s", err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get highlite items: response status is: %s", resp.Status())
	}

	reader := bytes.NewReader(resp.Body())

	// highlite export has Windows1257 encoding, need to convert it to utf-8
	decoder := transform.NewReader(reader, charmap.Windows1257.NewDecoder())

	return decoder, nil
}

// Gets highlite session cookies.
func (r *Client) getCookies(ctx context.Context) ([]*http.Cookie, error) {
	resp, err := r.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"Login":       r.login,
			"Password":    r.password,
			"LoginButton": "LoginButton",
		}).
		Post(r.loginEndpoint)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(resp.Status())
	}

	return resp.Cookies(), nil
}
