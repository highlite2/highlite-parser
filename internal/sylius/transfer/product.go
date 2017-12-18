package transfer

type Product struct {
	ID           int                    `json:"id,omitempty"`
	Code         string                 `json:"code,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
	Images       []Image                `json:"images,omitempty"`
}
