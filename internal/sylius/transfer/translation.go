package transfer

// Translation is a representation of translation in Sylius
type Translation struct {
	Description string `json:"description,omitempty"`
	ID          int    `json:"id,omitempty"`
	Locale      string `json:"locale,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
}
