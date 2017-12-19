package sylius

import (
	"context"

	"highlite-parser/internal/sylius/transfer"
)

// GetTaxon gets a category by code.
func (c *Client) GetTaxon(ctx context.Context, code string) (*transfer.Taxon, error) {
	result := &transfer.Taxon{}
	err := c.requestGet(ctx, c.getURL("/v1/taxons/%s", code), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateTaxon creates a taxon
func (c *Client) CreateTaxon(ctx context.Context, body transfer.TaxonNew) (*transfer.Taxon, error) {
	result := &transfer.Taxon{}
	err := c.requestPost(ctx, c.getURL("/v1/taxons/"), result, body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetProduct gets a product by code
func (c *Client) GetProduct(ctx context.Context, code string) (*transfer.Product, error) {
	result := &transfer.Product{}
	err := c.requestGet(ctx, c.getURL("/v1/products/%s", code), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateProduct creates a product.
func (c *Client) CreateProduct(ctx context.Context, body transfer.ProductNew) (*transfer.Product, error) {
	result := &transfer.Product{}
	err := c.requestPost(ctx, c.getURL("/v1/products/"), result, body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateProductVariant creates a product.
func (c *Client) CreateProductVariant(ctx context.Context, product string, body transfer.ProductVariantNew) (*transfer.ProductVariant, error) {
	result := &transfer.ProductVariant{}
	err := c.requestPost(ctx, c.getURL("/v1/products/%s/variants/", product), result, body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
