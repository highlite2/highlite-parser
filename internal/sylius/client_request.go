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

// TODO
func (c *Client) getRequestWithToken(ctx context.Context) (*resty.Request, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	request := resty.R().SetContext(ctx).SetHeader("Authorization", "Bearer "+token)

	return request, nil
}

// TODO
func (c *Client) executeRequestWithMethod(request *resty.Request, method string, url string) (*resty.Response, error) {
	switch method {
	case methodGet:
		return request.Get(url)
	case methodPost:
		return request.Post(url)
	case methodPatch:
		return request.Patch(url)
	case methodPut:
		return request.Put(url)
	}

	return nil, fmt.Errorf("unknown method")
}

// TODO
func (c *Client) checkResponseStatus(response *resty.Response) error {
	switch response.StatusCode() {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	}

	c.logger.Errorf("[%s] %s %s %s", response.Request.Method, response.Request.URL, response.Status(), string(response.Body()))

	return fmt.Errorf("%s", response.Status())
}

// Performs a request. Sets authorization token and handles errors.
// Creates context with timeout.
func (c *Client) request(ctx context.Context, method string, url string, result interface{}, body interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	request, err := c.getRequestWithToken(ctx)
	if err != nil {
		return err
	}

	if result != nil {
		request.SetResult(result)
	}

	if body != nil {
		request.SetBody(body)
	}

	c.logger.Debugf("Performing [%s] request to %s", method, url)

	response, err := c.executeRequestWithMethod(request, method, url)
	if err != nil {
		c.logger.Errorf(err.Error())

		return fmt.Errorf(err.Error())
	}

	return c.checkResponseStatus(response)
}
