package sylius

import (
	"context"

	"highlite-parser/internal/sylius/transfer"
)

// GetTaxon gets a category by code.
func (c *client) GetTaxon(ctx context.Context, code string) (*transfer.Taxon, error) {
	result := &transfer.TaxonRaw{}
	err := c.requestGet(ctx, c.getURL("/v1/taxons/%s", code), result)
	if err != nil {
		return nil, err
	}

	return transfer.ConvertRawTaxon(result), nil
}

// CreateTaxon creates a taxon
func (c *client) CreateTaxon(ctx context.Context, body *transfer.TaxonNew) (*transfer.Taxon, error) {
	result := &transfer.TaxonRaw{}
	err := c.requestPost(ctx, c.getURL("/v1/taxons/"), result, body)
	if err != nil {
		return nil, err
	}

	return transfer.ConvertRawTaxon(result), nil
}
