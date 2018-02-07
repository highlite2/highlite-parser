package form

import (
	"testing"

	"highlite-parser/internal/sylius/transfer"

	"github.com/stretchr/testify/assert"
)

func TestDisplay(t *testing.T) {
	// arrange
	enc := NewEncoder(getEncoderTestData())
	enc.Tag = "json"
	enc.PathToString = func(path []string) string {
		if len(path) > 0 && path[0] == "ProductEntire" {
			path = path[1:]
		}

		return PathToString(path)
	}

	// act
	values, err := enc.Values()

	// assert
	assert.Nil(t, err)
	assert.Equal(t, getEncoderExpectedOutput(), values)
}

func getEncoderTestData() transfer.Product {
	return transfer.Product{
		ProductEntire: transfer.ProductEntire{
			Code: "123123",
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
			Enabled: false,
		},

		MainTaxon:     "main taxon",
		ProductTaxons: "",
		Channels:      []string{"ch1", "ch2"},
	}
}

func getEncoderExpectedOutput() map[string]string {
	return map[string]string{
		"code":                                  "123123",
		"enabled":                               "",
		"mainTaxon":                             "main taxon",
		"channels[0]":                           "ch1",
		"channels[1]":                           "ch2",
		"translations[ru_RU][name]":             "name",
		"translations[ru_RU][slug]":             "slug",
		"translations[ru_RU][description]":      "description",
		"translations[ru_RU][shortDescription]": "short description",
	}
}
