package transfer

import (
	"encoding/json"
)

// Taxon is a representation of a category in Sylius.
type Taxon struct {
	ID           int                    `json:"id,omitempty"`
	Code         string                 `json:"code,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Parent       *Taxon                 `json:"parent,omitempty"`
	Root         *Taxon                 `json:"root,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
}

// TaxonRaw is a helper to parse Sylius taxon.
// Sylius api returns different representation for Translations:
// if there are no translations, Sylius returns an empty array,
// if not empty - an object.
type TaxonRaw struct {
	Taxon
	Parent       *TaxonRaw       `json:"parent"`
	Root         *TaxonRaw       `json:"root"`
	Translations json.RawMessage `json:"translations"`
}

// ConvertRawTaxon converts TaxonRaw to Taxon type.
func ConvertRawTaxon(raw *TaxonRaw) *Taxon {
	taxon := raw.Taxon

	if raw.Root != nil {
		taxon.Root = ConvertRawTaxon(raw.Root)
	}

	if raw.Parent != nil {
		taxon.Parent = ConvertRawTaxon(raw.Parent)
	}

	var t map[string]Translation
	if err := json.Unmarshal(raw.Translations, &t); err == nil {
		taxon.Translations = t
	}

	return &taxon
}
