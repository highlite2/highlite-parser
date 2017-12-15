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
	if _, err := w.importCategoryWithMemo(ctx, p.Category1); err != nil {
		return err
	}

	if _, err := w.importCategoryWithMemo(ctx, p.Category2); err != nil {
		return err
	}

	if _, err := w.importCategoryWithMemo(ctx, p.Category3); err != nil {
		return err
	}

	return nil
}

// Imports a category. Is a wrapper for importCategory method. Stores the result in local memory.
// Concurrent requests for the same key are blocked until the first completes.
func (w *Importer) importCategoryWithMemo(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	data, err := w.memo.Get(category.GetCode(), func() (interface{}, error) {
		return w.importCategory(ctx, category)
	})

	if err != nil {
		return nil, err
	}

	taxon, ok := data.(*transfer.Taxon)
	if !ok {
		return nil, fmt.Errorf("can't cast to *transfer.Taxon: %#v", data)
	}

	return taxon, nil
}

// Imports a category. First it tries to get the category from sylius api and if it doesn't - creates a new one.
func (w *Importer) importCategory(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	taxon, err := w.client.GetTaxon(ctx, category.GetCode())
	if err == nil {
		return taxon, nil
	}

	if err == sylius.ErrNotFound {
		data := CreateNewTaxonFromHighliteCategory(category)
		taxon, err := w.client.CreateTaxon(ctx, data)
		if err != nil {
			return nil, err
		}

		return taxon, nil
	}

	return nil, err
}
