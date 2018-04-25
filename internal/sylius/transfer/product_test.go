package transfer

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductEntire_UnmarshalJSON(t *testing.T) {
	// act
	product := getProductFromJSON()

	// assert
	assert.Equal(t, "highlite-40060", product.Code)
	assert.Equal(t, "entertainment_lighting_moving_heads_moving_heads_panels", product.MainTaxon.Code)
	assert.Equal(t, 3, len(product.ProductTaxons))
	assert.Equal(t, "entertainment_lighting", product.ProductTaxons[2].Taxon.Code)
	assert.Equal(t, "entertainment_lighting_moving_heads", product.ProductTaxons[1].Taxon.Code)
	assert.Equal(t, "entertainment_lighting_moving_heads_moving_heads_panels", product.ProductTaxons[0].Taxon.Code)
	assert.Equal(t, 20, len(product.Images))
	assert.Equal(t, "phantom-300-led-matrix-40060", product.Translations["en_US"].Slug)
	assert.Equal(t, "phantom-300-led-matrix-40060", product.Translations["ru_RU"].Slug)
	assert.Equal(t, 1, len(product.Channels))
	assert.Equal(t, "default", product.Channels[0].Code)
	assert.Equal(t, "mug_material", product.RawAttributes[0].Code)
	assert.Equal(t, "en_US", product.RawAttributes[0].LocaleCode)
	assert.Equal(t, "select", product.RawAttributes[0].Type)
	assert.Equal(t, "concrete", product.RawAttributes[0].Value)
}

func getProductFromJSON() ProductEntire {
	bytes, err := ioutil.ReadFile("_test_data/product-example.json")
	if err != nil {
		panic(err)
	}

	product := ProductEntire{}
	err = json.Unmarshal(bytes, &product)
	if err != nil {
		panic(err)
	}

	return product
}
