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

// NewImporter creates new Importer.
func NewImporter(client sylius.IClient, memo cache.IMemo, logger log.ILogger) *Importer {
	return &Importer{
		client: client,
		memo:   memo,
		logger: logger,
	}
}

// Importer imports highlite product into sylius.
type Importer struct {
	client sylius.IClient
	memo   cache.IMemo
	logger log.ILogger
}

// ImportProduct imports highlite product into sylius.
func (w *Importer) ImportProduct(ctx context.Context, p highlite.Product) error {
	_, err := w.importCategory(ctx, p.Category3)
	if err != nil {
		return err
	}

	return nil
}

// Imports a category. First tries to check if the category exists. It it doesn't, it makes
// a recursive call to import a parent category (if there is a link to parent category). After
// parent category is imported (actually it could be cached already) - it creates the current
// category.
func (w *Importer) importCategory(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	taxon, err := w.memoGetCategory(ctx, category)
	if err == nil {
		return taxon, nil

	} else if err == sylius.ErrNotFound {
		if category.Parent != nil {
			_, err := w.importCategory(ctx, category.Parent)
			if err != nil {
				return nil, err
			}
		}

		_, err := w.memoCreateCategory(ctx, category)
		if err != nil {
			return nil, err
		}

		return taxon, nil

	} else {
		return nil, err
	}
}

// Tries to find a category. Stores the result in local memory. Concurrent
// requests for the same key are blocked until the first completes.
func (w *Importer) memoGetCategory(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	data, err := w.memo.Get(category.GetCode(), func() (interface{}, error) {
		return w.client.GetTaxon(ctx, category.GetCode())
	})

	if err != nil {
		return nil, err
	}

	return castToTaxon(data)
}

// Imports a category. Stores the result in local memory. Concurrent
// requests for the same key are blocked until the first completes.
func (w *Importer) memoCreateCategory(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	data, err := w.memo.Get(category.GetCode(), func() (interface{}, error) {
		return w.client.CreateTaxon(ctx, CreateNewTaxonFromHighliteCategory(category))
	})

	if err != nil {
		return nil, err
	}

	return castToTaxon(data)
}

// Casts an interface{} type to *transfer.Taxon
func castToTaxon(data interface{}) (*transfer.Taxon, error) {
	taxon, ok := data.(*transfer.Taxon)
	if !ok {
		return nil, fmt.Errorf("can't cast to *transfer.Taxon: %#v", data)
	}

	return taxon, nil
}
