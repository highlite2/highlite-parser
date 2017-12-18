package internal

import (
	"context"
	"fmt"

	"highlite-parser/internal/cache"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"
	"highlite-parser/internal/sylius/transfer"
)

// NewCategoryImport creates new CategoryImport.
func NewCategoryImport(client sylius.IClient, memo cache.IMemo, logger log.ILogger) *CategoryImport {
	return &CategoryImport{
		client: client,
		memo:   memo,
		logger: logger,
	}
}

// CategoryImport imports highlite product into sylius.
type CategoryImport struct {
	client sylius.IClient
	memo   cache.IMemo
	logger log.ILogger
}

// Import imports a categoryImport. First tries to check if the categoryImport exists. If it doesn't, makes
// a recursive call to import a parent categoryImport (if there is a link to parent categoryImport). After
// parent categoryImport is imported (actually it could be cached already) - it creates the current
// categoryImport.
func (i *CategoryImport) Import(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	taxon, err := i.memoGetCategory(ctx, category)
	if err == nil {
		return taxon, nil

	} else if err == sylius.ErrNotFound {
		if category.Parent != nil {
			_, err := i.Import(ctx, category.Parent)
			if err != nil {
				return nil, err
			}
		}

		_, err := i.memoCreateCategory(ctx, category)
		if err != nil {
			return nil, err
		}

		return taxon, nil

	} else {
		return nil, err
	}
}

// Tries to find a categoryImport. Stores the result in local memory. Concurrent
// requests for the same key are blocked until the first completes.
func (i *CategoryImport) memoGetCategory(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	data, err := i.memo.Get(category.GetCode(), func() (interface{}, error) {
		return i.client.GetTaxon(ctx, category.GetCode())
	})

	if err != nil {
		return nil, err
	}

	return castInterfaceToTaxon(data)
}

// Imports a categoryImport. Stores the result in local memory. Concurrent
// requests for the same key are blocked until the first completes.
func (i *CategoryImport) memoCreateCategory(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	data, err := i.memo.Get(category.GetCode(), func() (interface{}, error) {
		return i.client.CreateTaxon(ctx, i.createNewTaxonFromHighliteCategory(category))
	})

	if err != nil {
		return nil, err
	}

	return castInterfaceToTaxon(data)
}

// Casts an interface{} type to *transfer.Taxon.
func castInterfaceToTaxon(data interface{}) (*transfer.Taxon, error) {
	taxon, ok := data.(*transfer.Taxon)
	if !ok {
		return nil, fmt.Errorf("can't cast to *transfer.Taxon: %#v", data)
	}

	return taxon, nil
}

// Converts highlite categoryImport to sylius taxon struct.
func (i *CategoryImport) createNewTaxonFromHighliteCategory(cat *highlite.Category) transfer.TaxonNew {
	taxon := transfer.TaxonNew{
		Code: cat.GetCode(),
		Translations: map[string]transfer.Translation{ // TODO take info from available locales from config
			transfer.LocaleEn: {
				Name: cat.Name,
				Slug: cat.GetURL(),
			},
			transfer.LocaleRu: {
				Name: cat.Name,
				Slug: cat.GetURL(),
			},
		},
	}

	if cat.Parent != nil {
		taxon.Parent = cat.Parent.GetCode()
	}

	return taxon
}
