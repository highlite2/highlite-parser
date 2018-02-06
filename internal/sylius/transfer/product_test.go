package transfer

import (
	"fmt"
	"testing"

	"github.com/fatih/structs"
)

func TestConvert(t *testing.T) {
	pr := getProduct()
	dst := structs.Map(pr)

	fmt.Println(dst)
}

func getProduct() Product {
	return Product{
		ProductEntire: ProductEntire{
			Code: "code 123",
			Images: []Image{
				{
					Type: "image type",
					Path: "image path",
				},
			},
			Translations: map[string]Translation{
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
