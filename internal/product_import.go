package internal

import (
	"context"
	"strings"

	"highlite-parser/internal/cache"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"
	"highlite-parser/internal/sylius/transfer"
)

// NewProductImport creates new ProductImport.
func NewProductImport(client sylius.IClient, memo cache.IMemo, logger log.ILogger) *ProductImport {
	return &ProductImport{
		client:         client,
		categoryImport: NewCategoryImport(client, memo, logger),
	}
}

// ProductImport imports highlite product into sylius.
type ProductImport struct {
	client         sylius.IClient
	categoryImport *CategoryImport
}

// Import imports highlite product into sylius.
func (i *ProductImport) Import(ctx context.Context, high highlite.Product) error {
	if product, err := i.client.GetProduct(ctx, high.Code); err == nil {
		return i.updateProduct(ctx, product, high)
	} else if err == sylius.ErrNotFound {
		return i.createProduct(ctx, high)
	} else {
		return err
	}
}

// Creates product.
func (i *ProductImport) createProduct(ctx context.Context, high highlite.Product) error {
	if _, err := i.categoryImport.Import(ctx, high.Category3); err != nil {
		return err
	}

	product := i.getProductFromHighlite(transfer.ProductEntire{}, high)
	productEntire, err := i.client.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	variant := i.getVariantFromHighlite(transfer.VariantEntire{}, high)
	if _, err := i.client.CreateProductVariant(ctx, productEntire.Code, variant); err != nil {
		return err
	}

	return nil
}

// Updates product.
func (i *ProductImport) updateProduct(ctx context.Context, product *transfer.ProductEntire, high highlite.Product) error {
	if err := i.client.UpdateProduct(ctx, i.getProductFromHighlite(*product, high)); err != nil {
		return err
	}

	variantCode := getProductMainVariantCode(product.Code)
	if variantEntire, err := i.client.GetProductVariant(ctx, product.Code, variantCode); err != nil {
		if err != sylius.ErrNotFound {
			return err
		}

		variant := i.getVariantFromHighlite(transfer.VariantEntire{}, high)
		if _, err := i.client.CreateProductVariant(ctx, product.Code, variant); err != nil {
			return err
		}
	} else {
		variant := i.getVariantFromHighlite(*variantEntire, high)
		if err := i.client.UpdateProductVariant(ctx, product.Code, variant); err != nil {
			return err
		}
	}

	return nil
}

// Creates sylius Variant structure from higlite product structure.
func (i *ProductImport) getVariantFromHighlite(variantEntire transfer.VariantEntire, high highlite.Product) transfer.Variant {
	channel := "US_WEB" // TODO take it from config

	variant := transfer.Variant{VariantEntire: variantEntire}

	variant.Code = getProductMainVariantCode(high.Code)
	variant.Translations = map[string]transfer.Translation{
		transfer.LocaleEn: {
			Name: high.Name,
		},
		transfer.LocaleRu: {
			Name: high.Name,
		},
	}
	variant.ChannelPrices = map[string]transfer.ChannelPrice{
		channel: {
			Price: high.Price,
		},
	}
	variant.Width = high.Width
	variant.Height = high.Height
	variant.Weight = high.Weight
	variant.Depth = high.Length

	return variant
}

// Creates sylius Product structure from higlite product structure.
func (i *ProductImport) getProductFromHighlite(productEntire transfer.ProductEntire, high highlite.Product) transfer.Product {
	channel := "US_WEB" // TODO take it from config

	product := transfer.Product{ProductEntire: productEntire}
	product.Code = high.Code
	product.MainTaxon = high.Category3.GetCode()
	product.ProductTaxons = strings.Join([]string{high.Category3.GetCode(), high.Category2.GetCode(), high.Category1.GetCode()}, ",")
	product.Channels = []string{channel}
	product.Translations = map[string]transfer.Translation{
		transfer.LocaleEn: {
			Name:             high.Name,
			Slug:             high.URL,
			Description:      high.ProductDescription(),
			ShortDescription: high.SubHeading,
		},
		transfer.LocaleRu: { // TODO "temporary" don't overwrite Russian description - it should be taken from Russian translations
			Name:             high.Name,
			Slug:             high.URL,
			Description:      high.ProductDescription(),
			ShortDescription: high.SubHeading,
		},
	}
	product.Enabled = true

	switch high.Status {
	case highlite.StatusDecline, highlite.StatusEOL:
		product.Enabled = false
	}

	return product
}

func getProductMainVariantCode(productCode string) string {
	return productCode + "_main"
}
