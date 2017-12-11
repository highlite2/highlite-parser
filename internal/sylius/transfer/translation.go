package transfer

// Translation is a representation of translation in Sylius
type Translation struct {
	Description string `json:"description"`
	ID          int    `json:"id"`
	Locale      string `json:"locale"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
}
