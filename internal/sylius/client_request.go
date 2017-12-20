package sylius

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty"
)

// Performs GET request.
func (c *Client) requestGet(ctx context.Context, url string, result interface{}) error {
	return c.request(ctx, methodGet, url, result, nil)
}

// Performs POST request.
func (c *Client) requestPost(ctx context.Context, url string, result interface{}, body interface{}) error {
	return c.request(ctx, methodPost, url, result, body)
}

// Performs PATCH request.
func (c *Client) requestPatch(ctx context.Context, url string, body interface{}) error {
	return c.request(ctx, methodPatch, url, nil, body)
}

// Performs a request. Sets authorization token and handles errors.
// Creates context with timeout.
// TODO this method seems to be long, should consider to split it.
func (c *Client) request(ctx context.Context, method string, url string, result interface{}, body interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	var err error

	token, err := c.getToken()
	if err != nil {
		c.logger.Errorf("Failed to get token during [%s] %s request", method, url)

		return err
	}

	request := resty.R().SetContext(ctx).SetHeader("Authorization", "Bearer "+token)
	if result != nil {
		request.SetResult(result)
	}
	if body != nil {
		request.SetBody(body)
	}

	var res *resty.Response

	c.logger.Debugf("Performing [%s] request to %s", method, url)

	switch method {
	case methodGet:
		res, err = request.Get(url)
	case methodPost:
		res, err = request.Post(url)
	case methodPatch:
		res, err = request.Patch(url)
	default:
		err = fmt.Errorf("unknown method")
	}

	if err != nil {
		c.logger.Errorf(err.Error())

		return fmt.Errorf(err.Error())
	}

	switch res.StatusCode() {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	}

	c.logger.Errorf("[%s] %s %s %s", method, url, res.Status(), string(res.Body()))

	return fmt.Errorf("%s", res.Status())
}
