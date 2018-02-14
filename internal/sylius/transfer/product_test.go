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
}

func TestProductsEqual_Equal(t *testing.T) {
	// arrange
	entire := getProductFromJSON()
	product := Product{
		ProductEntire: entire,
		MainTaxon:     "entertainment_lighting_moving_heads_moving_heads_panels",
		ProductTaxons: "entertainment_lighting, entertainment_lighting_moving_heads, entertainment_lighting_moving_heads_moving_heads_panels",
		Channels:      []string{"default"},
	}

	// act
	// assert
	assert.True(t, ProductsEqual(entire, product))
}

func TestProductsEqual_NotEqual_1(t *testing.T) {
	// arrange
	product := getEqualProduct()
	product.Code = "wrong"

	// act
	// assert
	assert.False(t, ProductsEqual(getProductFromJSON(), product))
}

func getEqualProduct() Product {
	return Product{
		ProductEntire: getProductFromJSON(),
		MainTaxon:     "entertainment_lighting_moving_heads_moving_heads_panels",
		ProductTaxons: "entertainment_lighting, entertainment_lighting_moving_heads, entertainment_lighting_moving_heads_moving_heads_panels",
		Channels:      []string{"default"},
	}
}

var getProductFromJSONCache *ProductEntire

func getProductFromJSON() ProductEntire {
	if getProductFromJSONCache != nil {
		return *getProductFromJSONCache
	}

	bytes, err := ioutil.ReadFile("_test_data/product-example.json")
	if err != nil {
		panic(err)
	}

	product := ProductEntire{}
	err = json.Unmarshal(bytes, &product)
	if err != nil {
		panic(err)
	}

	getProductFromJSONCache = &product

	return product
}
