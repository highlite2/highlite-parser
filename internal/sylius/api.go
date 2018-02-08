package sylius

import (
	"context"
	"fmt"
	"time"

	"highlite-parser/internal/form"
	"highlite-parser/internal/sylius/transfer"
)

// GetTaxon gets a category by code.
func (c *Client) GetTaxon(ctx context.Context, taxon string) (*transfer.Taxon, error) {
	result := &transfer.Taxon{}
	err := c.requestGet(ctx, c.getURL("/v1/taxons/%s", taxon), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateTaxon creates a taxon.
func (c *Client) CreateTaxon(ctx context.Context, taxon transfer.TaxonNew) (*transfer.Taxon, error) {
	result := &transfer.Taxon{}
	err := c.requestPost(ctx, c.getURL("/v1/taxons/"), result, taxon)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetProduct gets a product by code.
func (c *Client) GetProduct(ctx context.Context, product string) (*transfer.ProductEntire, error) {
	result := &transfer.ProductEntire{}
	err := c.requestGet(ctx, c.getURL("/v1/products/%s", product), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateProduct creates a product.
func (c *Client) CreateProduct(ctx context.Context, product transfer.Product, images []transfer.ImageUpload) (*transfer.ProductEntire, error) {
	url := c.getURL("/v1/products/")

	defer c.timeTrack(time.Now(), fmt.Sprintf("Performing [%s] request to %s", methodPost, url))

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	result := &transfer.ProductEntire{}

	request, err := c.getRequestWithToken(ctx)
	if err != nil {
		return nil, err
	}

	formEncoder := form.NewEncoder(product)
	formEncoder.FieldTag = "json"
	values, err := formEncoder.Values()
	if err != nil {
		return nil, err
	}

	request.SetFormData(values)
	request.SetResult(result)

	for i, image := range images {
		request.SetFileReader(fmt.Sprintf("images[%d][file]", i), image.Name, image.Reader)
	}

	if response, err := c.executeRequestWithMethod(request, methodPost, url); err != nil {
		c.logger.Errorf(err.Error())

		return nil, err
	} else if err := c.checkResponseStatus(response); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateProduct updates product.
func (c *Client) UpdateProduct(ctx context.Context, product transfer.Product) error {
	return c.requestPatch(ctx, c.getURL("/v1/products/%s", product.Code), product)
}

// GetProductVariant gets a variant by product code and variant code.
func (c *Client) GetProductVariant(ctx context.Context, product string, variant string) (*transfer.VariantEntire, error) {
	result := &transfer.VariantEntire{}
	url := c.getURL("/v1/products/%s/variants/%s", product, variant)
	if err := c.requestGet(ctx, url, result); err != nil {
		return result, err
	}

	return result, nil
}

// CreateProductVariant creates a product.
func (c *Client) CreateProductVariant(ctx context.Context, product string, variant transfer.Variant) (*transfer.VariantEntire, error) {
	result := &transfer.VariantEntire{}
	err := c.requestPost(ctx, c.getURL("/v1/products/%s/variants/", product), result, variant)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateProductVariant creates a product.
func (c *Client) UpdateProductVariant(ctx context.Context, product string, variant transfer.Variant) error {
	url := c.getURL("/v1/products/%s/variants/%s", product, variant.Code)

	return c.requestPatch(ctx, url, variant)
}
