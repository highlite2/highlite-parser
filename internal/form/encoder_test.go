package form

import (
	"testing"
	"highlite-parser/internal/sylius/transfer"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestDisplay(t *testing.T) {
	enc := NewEncoder(getProduct())
	enc.tag = "json"

	values, err := enc.Values()

	assert.Nil(t, err)

	for key, val := range values {
		fmt.Printf("%s: %s\n", key, val)
	}
}

func getProduct() transfer.Product {
	return transfer.Product{
		ProductEntire: transfer.ProductEntire{
			Code: "code 123",
			Images: []transfer.Image{
				{
					Type: "image type",
					Path: "image path",
				},
			},
			Translations: map[string]transfer.Translation{
				"ru_RU": {
					Name:             "name",
					Slug:             "slug",
					Description:      "description",
					ShortDescription: "short description",
				},
			},
			Enabled: true,
		},

		MainTaxon:     "main taxon",
		ProductTaxons: "taxon 1, taxon 2",
		Channels:      []string{"ch1", "ch2"},
	}
}