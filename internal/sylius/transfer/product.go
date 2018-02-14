package transfer

import (
	"encoding/json"
	"strings"
)

// ProductEntire is a representation of a product in Sylius.
type ProductEntire struct {
	Code          string                 `json:"code,omitempty"`
	Translations  map[string]Translation `json:"translations,omitempty"`
	Images        []Image                `json:"images,omitempty"`
	Enabled       bool                   `json:"enabled"`
	ProductTaxons []TaxonWrap            `json:"productTaxons,omitempty"`
	MainTaxon     *Taxon                 `json:"mainTaxon,omitempty"`
	Channels      []Channel              `json:"channels"`
}

// Help structure to unmarshal Sylius api response.
type productEntireRaw struct {
	Code          string                 `json:"code,omitempty"`
	Translations  map[string]Translation `json:"translations,omitempty"`
	Images        json.RawMessage        `json:"images,omitempty"`
	Enabled       bool                   `json:"enabled"`
	ProductTaxons []TaxonWrap            `json:"productTaxons,omitempty"`
	MainTaxon     *Taxon                 `json:"mainTaxon,omitempty"`
	Channels      []Channel              `json:"channels"`
}

// UnmarshalJSON helps to fix inconsistency in sylius api response.
// Sylius returns image as a slice or, sometimes, as a map.
func (p *ProductEntire) UnmarshalJSON(value []byte) error {
	raw := &productEntireRaw{}
	if err := json.Unmarshal(value, raw); err != nil {
		return err
	}

	p.Code = raw.Code
	p.Translations = raw.Translations
	p.Enabled = raw.Enabled
	p.MainTaxon = raw.MainTaxon
	p.ProductTaxons = raw.ProductTaxons
	p.Channels = raw.Channels

	var images []Image
	if err := json.Unmarshal(raw.Images, &images); err == nil {
		p.Images = images
	} else {
		var imageMap map[string]Image
		if err := json.Unmarshal(raw.Images, &imageMap); err == nil {
			images = make([]Image, len(imageMap))
			i := 0
			for _, im := range imageMap {
				images[i] = im
				i++
			}
			p.Images = images
		}
	}

	return nil
}

// Product is a structure to be used in product create/update requests.
type Product struct {
	ProductEntire

	MainTaxon     string   `json:"mainTaxon,omitempty"`
	ProductTaxons string   `json:"productTaxons,omitempty"` // String in which the codes of taxons was written down (separated by comma)
	Channels      []string `json:"channels,omitempty"`
}

// ProductsEqual checks if api response product equals to the composed one.
// It doesn't check Images and Enabled flag.
func ProductsEqual(e ProductEntire, p Product) bool {
	if e.Code != p.Code {
		return false
	}

	// checking main taxon
	if p.MainTaxon != "" {
		if e.MainTaxon == nil || p.MainTaxon != e.MainTaxon.Code {
			return false
		}
	} else if e.MainTaxon != nil {
		return false
	}

	// checking taxons
	taxons := strings.Split(p.ProductTaxons, ",")
	if len(taxons) != len(e.ProductTaxons) {
		return false
	}
	for _, taxon := range taxons {
		taxon = strings.Trim(taxon, " ")
		found := false
		for _, wrap := range e.ProductTaxons {
			if wrap.Taxon.Code == taxon {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// checking translations
	if len(e.Translations) != len(p.Translations) {
		return false
	}
	for k, ptr := range p.Translations {
		etr, ok := e.Translations[k]
		if !ok || ptr.Name != etr.Name || ptr.Description != etr.Description ||
			ptr.Slug != etr.Slug || ptr.ShortDescription != etr.ShortDescription {
			return false
		}
	}

	// checking channels
	if len(p.Channels) != len(e.Channels) {
		return false
	}
	for _, code := range p.Channels {
		found := false
		for _, channel := range e.Channels {
			if channel.Code == code {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
