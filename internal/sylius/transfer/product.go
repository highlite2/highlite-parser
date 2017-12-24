package transfer

// ProductEntire is a representation of a product in Sylius.
type ProductEntire struct {
	Code         string                 `json:"code,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
	Images       []Image                `json:"-"`
	Enabled      bool                   `json:"enabled"`
}

// Product is a structure to be used in product create/update requesta.
type Product struct {
	ProductEntire

	MainTaxon     string   `json:"mainTaxon,omitempty"`
	ProductTaxons string   `json:"productTaxons,omitempty"` // String in which the codes of taxons was written down (separated by comma)
	Channels      []string `json:"channels,omitempty"`
}
