package transfer

// Taxon is a representation of a category in Sylius
type Taxon struct {
	ID           int                    `json:"id"`
	Code         string                 `json:"code"`
	Name         string                 `json:"name"`
	Parent       *Taxon                 `json:"parent"`
	Root         *Taxon                 `json:"root"`
	Translations map[string]Translation `json:"translations"`
}
