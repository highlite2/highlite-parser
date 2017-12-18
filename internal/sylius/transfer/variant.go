package transfer

// Variant is a representation of a product variant in Sylius.
type ProductVariant struct {
	ID           int                    `json:"id,omitempty"`
	Code         string                 `json:"code,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
}

// ProductVariantNew is a structure to be used in new product variant request.
type ProductVariantNew struct {
	ID           int                    `json:"id,omitempty"`
	Code         string                 `json:"code,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
	OnHand       int                    `json:"onHand,omitempty"`
	Width        float64                `json:"width,omitempty"`
	Height       float64                `json:"height,omitempty"`
	Depth        float64                `json:"depth,omitempty"`
	Weight       float64                `json:"weight,omitempty"`
}
