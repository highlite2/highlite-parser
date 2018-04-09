package transfer

import (
	"encoding/json"
	"fmt"
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
	assert.Equal(t, "mug_color", product.Attributes[0].Attribute)
	assert.Equal(t, "en_US", product.Attributes[0].LocaleCode)
	assert.Equal(t, "yellow", product.Attributes[0].Value)
}

func TestProductsEqual_Equal(t *testing.T) {
	for i, test := range testProductsEqualCases {
		// arrange
		filter := test.filter(getProductEntireEqualProductMock())

		// act
		// assert
		assert.Equal(t, test.updateRequired, ProductUpdateRequired(getProductEntireMock(), filter), fmt.Sprintf("Failed test with index %d", i))
	}
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

func getProductEntireEqualProductMock() Product {
	return Product{
		ProductEntire: getProductEntireMock(),
		MainTaxon:     "taxon1",
		ProductTaxons: "taxon1, taxon2, taxon3",
		Channels:      []string{"channel1", "channel2"},
	}
}

func getProductEntireMock() ProductEntire {
	return ProductEntire{
		Code: "code",
		Translations: map[string]Translation{
			"en_US": {
				Name:             "name",
				Description:      "description",
				Slug:             "slug",
				ShortDescription: "short description",
			},
		},
		MainTaxon: &Taxon{
			Code: "taxon1",
		},
		ProductTaxons: []TaxonWrap{
			{Taxon: Taxon{Code: "taxon1"}},
			{Taxon: Taxon{Code: "taxon2"}},
			{Taxon: Taxon{Code: "taxon3"}},
		},
		Channels: []Channel{
			{Code: "channel1"},
			{Code: "channel2"},
		},
		Attributes: []ProductAttribute{
			{
				Attribute:  "type1",
				LocaleCode: "locale1",
				Value:      "value1",
			},
		},
	}
}

var testProductsEqualCases = []struct {
	updateRequired bool
	filter         func(v Product) Product
}{
	{
		updateRequired: true,
		filter: func(v Product) Product {
			v.Attributes = []ProductAttribute{}
			return v
		},
	},
	{
		updateRequired: true,
		filter: func(v Product) Product {
			v.Attributes = []ProductAttribute{
				{
					Attribute:  "type1",
					LocaleCode: "locale1",
					Value:      "value1",
				},
				{
					Attribute:  "type2",
					LocaleCode: "locale1",
					Value:      "value1",
				},
			}
			return v
		},
	},
	{
		filter: func(v Product) Product {
			v.Attributes = []ProductAttribute{
				{
					Attribute:  "type1",
					LocaleCode: "locale1",
					Value:      "value1",
				},
			}
			return v
		},
	},
	{
		updateRequired: true,
		filter: func(v Product) Product {
			v.Attributes = []ProductAttribute{
				{
					Attribute:  "type1",
					LocaleCode: "locale1",
					Value:      "value1",
				},
				{
					Attribute:  "type1",
					LocaleCode: "locale2",
					Value:      "value1",
				},
			}
			return v
		},
	},
	{
		updateRequired: true,
		filter: func(v Product) Product {
			v.Attributes = []ProductAttribute{
				{
					Attribute:  "type1",
					LocaleCode: "locale1",
					Value:      "value1",
				},
				{
					Attribute:  "type1",
					LocaleCode: "locale1",
					Value:      "value2",
				},
			}
			return v
		},
	},
	{
		filter: func(v Product) Product {
			return v
		},
	},
	{
		updateRequired: true,
		filter: func(v Product) Product {
			v.Code = ""
			return v
		},
	},
	{
		filter: func(v Product) Product {
			tr := v.Translations["en_US"]
			tr.Name = ""
			v.Translations["en_US"] = tr
			return v
		},
	},
	{
		filter: func(v Product) Product {
			tr := v.Translations["en_US"]
			tr.Description = ""
			v.Translations["en_US"] = tr
			return v
		},
	},
	{
		filter: func(v Product) Product {
			tr := v.Translations["en_US"]
			tr.Slug = ""
			v.Translations["en_US"] = tr
			return v
		},
	},
	{
		filter: func(v Product) Product {
			tr := v.Translations["en_US"]
			tr.ShortDescription = ""
			v.Translations["en_US"] = tr
			return v
		},
	},
	{
		filter: func(v Product) Product {
			tr := v.Translations["en_US"]
			v.Translations["new"] = tr
			return v
		},
	},
	{
		filter: func(v Product) Product {
			v.Translations = map[string]Translation{}
			return v
		},
	},
	{
		updateRequired: true,
		filter: func(v Product) Product {
			v.Channels[0] = "not default"
			return v
		},
	},
	{
		filter: func(v Product) Product {
			v.Channels[0], v.Channels[1] = v.Channels[1], v.Channels[0]
			return v
		},
	},
	{
		updateRequired: true,
		filter: func(v Product) Product {
			v.ProductTaxons = "taxon3, taxon1"
			return v
		},
	},
	{
		updateRequired: true,
		filter: func(v Product) Product {
			v.MainTaxon = ""
			return v
		},
	},
	{
		filter: func(v Product) Product {
			v.ProductTaxons = "taxon3, taxon2, taxon1"
			return v
		},
	},
}
