package transfer

import (
	"encoding/json"
)

// Taxon is a representation of a category in Sylius.
// IMPORTANT!!! If you add new fields here, you have to add the same fields
// to taxonRaw struct and extend UnmarshalJSON func (t *Taxon) UnmarshalJSON method.
type Taxon struct {
	ID           int                    `json:"id,omitempty"`
	Code         string                 `json:"code,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Parent       *Taxon                 `json:"parent,omitempty"`
	Root         *Taxon                 `json:"root,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
}

// A full copy of Taxon struct but has different Translations field definition.
// This struct helps to unmarshal json to Taxon. Sylius api returns different
// representation for Translations: if there are no translations, Sylius returns an
// empty array, if not empty - an object.
type taxonRaw struct {
	ID           int             `json:"id,omitempty"`
	Code         string          `json:"code,omitempty"`
	Name         string          `json:"name,omitempty"`
	Parent       *Taxon          `json:"parent,omitempty"`
	Root         *Taxon          `json:"root,omitempty"`
	Translations json.RawMessage `json:"translations"`
}

// TaxonWrap is used in sylius product response to wrap product taxons.
type TaxonWrap struct {
	ID    int   `json:"id"`
	Taxon Taxon `json:"taxon"`
}

// UnmarshalJSON helps to fix inconsistency in sylius api response.
func (t *Taxon) UnmarshalJSON(value []byte) error {
	raw := &taxonRaw{}
	if err := json.Unmarshal(value, raw); err != nil {
		return err
	}

	t.ID = raw.ID
	t.Code = raw.Code
	t.Name = raw.Name
	t.Parent = raw.Parent
	t.Root = raw.Root

	var tr map[string]Translation
	if err := json.Unmarshal(raw.Translations, &tr); err == nil {
		t.Translations = tr
	}

	return nil
}

// TaxonNew is a structure to be used in new taxon request.
type TaxonNew struct {
	Code         string                 `json:"code"`
	Parent       string                 `json:"parent,omitempty"`
	Translations map[string]Translation `json:"translations"`
}
