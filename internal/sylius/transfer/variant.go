package transfer

import "reflect"

// VariantEntire is a representation of a product variant in Sylius.
type VariantEntire struct {
	Code          string                  `json:"code,omitempty"`
	Translations  map[string]Translation  `json:"translations,omitempty"`
	ChannelPrices map[string]ChannelPrice `json:"channelPricings,omitempty"`
}

// Variant is a structure to be used in new product variant request.
type Variant struct {
	VariantEntire
}

// VariantsEqual checks if variants are equal
func VariantsEqual(e VariantEntire, v Variant) bool {
	if e.Code != v.Code {
		return false
	}

	if !reflect.DeepEqual(e.Translations, v.Translations) {
		return false
	}

	if !reflect.DeepEqual(e.ChannelPrices, v.ChannelPrices) {
		return false
	}

	return true
}
