package imprt

import (
	"context"
	"fmt"

	"highlite2-import/internal/cache"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/log"
	"highlite2-import/internal/sylius"
	"highlite2-import/internal/sylius/transfer"
)

// ICategoryImport imports imports a category.
type ICategoryImport interface {
	Import(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error)
}

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

// Import imports a category. First tries to check if the categoryImport exists. If it doesn't, makes
// a recursive call to import a parent categoryImport (if there is a link to parent categoryImport). After
// parent categoryImport is imported (actually it could be cached already) - it creates the current
// categoryImport.
func (i *CategoryImport) Import(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	taxon, err := i.memoGetCategory(ctx, category)
	if err == nil {
		return taxon, nil
	}

	if category.Parent != nil {
		_, err := i.Import(ctx, category.Parent)
		if err != nil {
			return nil, fmt.Errorf("%s (%s) parent category: %s", category.Name, category.GetCode(), err)
		}
	}

	_, err = i.memoCreateCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%s (%s) memoCreateCategory: %s", category.Name, category.GetCode(), err)
	}

	return taxon, nil
}

// Tries to find a categoryImport. Stores the result in local memory. Concurrent
// requests for the same key are blocked until the first completes.
func (i *CategoryImport) memoGetCategory(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	data, err := i.memo.GetOnce(category.GetCode(), func() (interface{}, error) {
		i.logger.Debugf("GetTaxon %s", category.GetCode())

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
		i.logger.Debugf("CreateTaxon %s", category.GetCode())

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
		Translations: map[string]transfer.Translation{
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
