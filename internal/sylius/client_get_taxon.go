package sylius

import (
	"context"
	"fmt"
	"net/http"

	"highlite-parser/internal/sylius/transfer"

	"github.com/go-resty/resty"
)

// GetTaxon get a category by its code.
func (c *client) GetTaxon(ctx context.Context, code string) (*transfer.Taxon, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	url := c.getURL("/v1/taxons/%s", code)

	c.log.Debug("getting token")
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	result := &transfer.TaxonRaw{}
	resp, err := resty.R().SetContext(ctx).SetHeader("Authorization", "Bearer "+token).SetResult(result).Get(url)

	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %s", url, err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request to %s ended with status %s", url, resp.Status())
	}

	return transfer.ConvertRawTaxon(result), nil
}
