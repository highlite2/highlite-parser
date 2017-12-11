package transfer

// Taxon is a representation of a category in Sylius
type Taxon struct {
	ID           int         `json:"id"`
	Code         string      `json:"code"`
	Name         string      `json:"name"`
	Parent       *Taxon      `json:"parent"`
	Root         *Taxon      `json:"root"`
	Translations interface{} `json:"translations"`
}

// GetTranslations returns translations map
func (t *Taxon) GetTranslations() map[string]Translation {
	var empty map[string]Translation

	if tr, ok := t.Translations.(map[string]Translation); ok {
		return tr
	}

	return empty
}
