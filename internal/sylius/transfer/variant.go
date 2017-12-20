package transfer

// VariantEntire is a representation of a product variant in Sylius.
type VariantEntire struct {
	Code          string                  `json:"code,omitempty"`
	Translations  map[string]Translation  `json:"translations,omitempty"`
	ChannelPrices map[string]ChannelPrice `json:"channelPricings,omitempty"`
	Width         float64                 `json:"width,omitempty"`
	Height        float64                 `json:"height,omitempty"`
	Depth         float64                 `json:"depth,omitempty"`
	Weight        float64                 `json:"weight,omitempty"`
	OnHand        int64                   `json:"onHand,omitempty"`
}

// Variant is a structure to be used in new product variant request.
type Variant struct {
	VariantEntire
}
