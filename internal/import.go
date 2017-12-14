package internal

import (
	"highlite-parser/internal/cache"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/sylius"
)

func NewImport(client sylius.IClient, memo cache.IMemo) *Import {
	return &Import{
		client: client,
		memo:   memo,
	}
}

type Import struct {
	client sylius.IClient
	memo   cache.IMemo
}

func (i *Import) ImportProduct(product *highlite.Product) {

}

func (i *Import) createTaxon(category *highlite.Category) {

}
