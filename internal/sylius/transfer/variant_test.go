package transfer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariantsEqual(t *testing.T) {
	for i, test := range testVariantsEqualCases {
		// arrange
		variant := Variant{VariantEntire: getVariantMock()}
		variant.ChannelPrices = map[string]ChannelPrice{
			"default": {
				Price: 1001 / 100.,
			},
		}

		filter := test.filter(variant)

		// act
		// assert
		assert.Equal(t, test.equal, VariantsEqual(getVariantMock(), filter), fmt.Sprintf("Failed test with index %d", i))
	}
}

var testVariantsEqualCases = []struct {
	equal  bool
	filter func(v Variant) Variant
}{
	{
		equal: true,
		filter: func(v Variant) Variant {
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			v.Code = ""
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			tr := v.Translations["en_US"]
			tr.Name = ""
			v.Translations["en_US"] = tr
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			tr := v.Translations["en_US"]
			tr.Description = ""
			v.Translations["en_US"] = tr
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			tr := v.Translations["en_US"]
			tr.Slug = ""
			v.Translations["en_US"] = tr
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			tr := v.Translations["en_US"]
			tr.ShortDescription = ""
			v.Translations["en_US"] = tr
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			tr := v.Translations["en_US"]
			v.Translations["new"] = tr
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			v.Translations = map[string]Translation{
				"en_US": v.Translations["en_US"],
			}
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			pr := v.ChannelPrices["default"]
			pr.Price = 100.1
			v.ChannelPrices["default"] = pr
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			pr := v.ChannelPrices["default"]
			v.ChannelPrices["new"] = pr
			return v
		},
	},
	{
		filter: func(v Variant) Variant {
			v.ChannelPrices = map[string]ChannelPrice{}
			return v
		},
	},
}

func getVariantMock() VariantEntire {
	return VariantEntire{
		Code: "variant",
		Translations: map[string]Translation{
			"en_US": {
				Name:             "name",
				Description:      "description",
				Slug:             "slug",
				ShortDescription: "short description",
			},
			"ru_RU": {
				Name:             "name 2",
				Description:      "description 2",
				Slug:             "slug 2",
				ShortDescription: "short description 2",
			},
		},
		ChannelPrices: map[string]ChannelPrice{
			"default": {
				Price: 1001,
			},
		},
	}
}
