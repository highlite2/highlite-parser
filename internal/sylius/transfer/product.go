package transfer

// Product is a representation of a product in Sylius.
type Product struct {
	ID           int                    `json:"id,omitempty"`
	Code         string                 `json:"code,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
	Images       []Image                `json:"images,omitempty"`
}

// ProductNew is a structure to be used in new product request.
type ProductNew struct {
	ID            int                    `json:"id,omitempty"`
	Code          string                 `json:"code,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Translations  map[string]Translation `json:"translations,omitempty"`
	Images        []Image                `json:"images,omitempty"`
	MainTaxon     string                 `json:"mainTaxon,omitempty"`
	ProductTaxons string                 `json:"productTaxons,omitempty"` // String in which the codes of taxons was written down (separated by comma)
	Channels      []string               `json:"channels,omitempty"`
}
