package internal

import (
	"context"
	"fmt"

	"highlite-parser/internal/cache"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/sylius"
)

// NewWriter creates new Writer
func NewWriter(client sylius.IClient, memo cache.IMemo) *Writer {
	return &Writer{
		client: client,
		memo:   memo,
	}
}

// Writer imports highlite product into sylius
type Writer struct {
	client sylius.IClient
	memo   cache.IMemo
}

// WriteProduct imports highlite product into sylius
func (w *Writer) WriteProduct(ctx context.Context, product highlite.Product) {
	taxon, err := w.client.GetTaxon(ctx, product.Category1.GetCode())
	if err == sylius.ErrNotFound {
		tr := CreateNewTaxonFromHighliteCategory(product.Category1)
		taxon, err := w.client.CreateTaxon(context.Background(), tr)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("%#v", taxon)
		}

	} else if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%#v", taxon)
	}
}
