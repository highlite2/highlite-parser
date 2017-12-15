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

// NewWriter creates new Writer
func NewWriter(client sylius.IClient, memo cache.IMemo, logger log.ILogger) *Writer {
	return &Writer{
		client: client,
		memo:   memo,
		logger: logger,
	}
}

// Writer imports highlite product into sylius
type Writer struct {
	client sylius.IClient
	memo   cache.IMemo
	logger log.ILogger
}

// WriteProduct imports highlite product into sylius
func (w *Writer) WriteProduct(ctx context.Context, p highlite.Product) error {
	if _, err := w.createCategoryWithMemo(ctx, p.Category1); err != nil {
		return err
	}

	if _, err := w.createCategoryWithMemo(ctx, p.Category2); err != nil {
		return err
	}

	if _, err := w.createCategoryWithMemo(ctx, p.Category3); err != nil {
		return err
	}

	return nil
}

// Creates a category. Is a wrapper for createCategory method. Stores the result in local memory.
// Uses local cache. Concurrent requests for the same key block until the first completes.
func (w *Writer) createCategoryWithMemo(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
	data, err := w.memo.Get(category.GetCode(), func() (interface{}, error) {
		return w.createCategory(ctx, category)
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

// Creates a category. First it tries to get the category and if it doesn't - creates a new one.
func (w *Writer) createCategory(ctx context.Context, category *highlite.Category) (*transfer.Taxon, error) {
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
