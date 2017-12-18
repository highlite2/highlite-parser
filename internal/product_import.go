package internal

import (
	"context"

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
func (i *ProductImport) Import(ctx context.Context, p highlite.Product) error {
	if _, err := i.client.GetProduct(ctx, p.Code); err != nil {
		if err != sylius.ErrNotFound {
			return err
		}

		if _, err := i.categoryImport.Import(ctx, p.Category3); err != nil {
			return err
		}

		_, err := i.client.CreateProduct(ctx, i.createNewProductFromHighliteProduct(p))
		if err != nil {
			return err
		}
	}

	return nil
}

// Converts highlite product to sylius product struct.
func (i *ProductImport) createNewProductFromHighliteProduct(p highlite.Product) transfer.ProductNew {
	product := transfer.ProductNew{
		Code: p.Code,
		Translations: map[string]transfer.Translation{ // TODO take info from available locales from config
			transfer.LocaleEn: {
				Name: p.Name,
				Slug: p.URL,
			},
			transfer.LocaleRu: {
				Name: p.Name,
				Slug: p.URL,
			},
		},
		MainTaxon:     p.Category3.GetCode(),
		ProductTaxons: p.Category3.GetCode(),
		Channels: []string{
			"US_WEB", // TODO take it from config
		},
	}

	return product
}
