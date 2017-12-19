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
func (i *ProductImport) Import(ctx context.Context, h highlite.Product) error {
	if product, err := i.client.GetProduct(ctx, h.Code); err != nil {
		if err != sylius.ErrNotFound {
			return err
		}

		if _, err := i.categoryImport.Import(ctx, h.Category3); err != nil {
			return err
		}

		data := i.getProductFromHighlite(transfer.ProductEntire{}, h)
		if !data.Enabled {
			return nil
		}

		product, err := i.client.CreateProduct(ctx, data)
		if err != nil {
			return err
		}

		variant := i.createNewProductVariantFromHighliteProduct(h)
		if _, err := i.client.CreateProductVariant(ctx, product.Code, variant); err != nil {
			return err
		}
	} else {
		if err := i.client.UpdateProduct(ctx, i.getProductFromHighlite(*product, h)); err != nil {
			return err
		}
	}

	return nil
}

// Converts highlite product to sylius product variant struct.
func (i *ProductImport) createNewProductVariantFromHighliteProduct(p highlite.Product) transfer.ProductVariantNew {
	channel := "US_WEB" // TODO take it from config

	variant := transfer.ProductVariantNew{
		Code: p.Code + "_main",
		Translations: map[string]transfer.Translation{
			transfer.LocaleEn: {
				Name: p.Name,
			},
			transfer.LocaleRu: {
				Name: p.Name,
			},
		},
		ChannelPrices: map[string]transfer.ChannelPrice{
			channel: {
				Price: p.Price,
			},
		},
		Width:  p.Width,
		Height: p.Height,
		Weight: p.Weight,
		Depth:  p.Length,
	}

	return variant
}

// Creates sylius Product structure from higlite product structure.
func (i *ProductImport) getProductFromHighlite(entire transfer.ProductEntire, h highlite.Product) transfer.Product {
	channel := "US_WEB" // TODO take it from config

	p := transfer.Product{ProductEntire: entire}
	p.Code = h.Code
	p.MainTaxon = h.Category3.GetCode()
	p.ProductTaxons = strings.Join([]string{h.Category3.GetCode(), h.Category2.GetCode(), h.Category1.GetCode()}, ",")
	p.Channels = []string{channel}
	p.Translations = map[string]transfer.Translation{
		transfer.LocaleEn: {
			Name:        h.Name,
			Slug:        h.URL,
			Description: h.Description,
		},
		transfer.LocaleRu: { // TODO "temporary" don't overwrite Russian description - it should be taken from Russian translations
			Name:        h.Name,
			Slug:        h.URL,
			Description: h.Description,
		},
	}
	p.Enabled = true

	switch h.Status {
	case highlite.StatusDecline, highlite.StatusEOL:
		p.Enabled = false
	}

	return p
}
