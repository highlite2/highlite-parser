package transfer

const (
	// LocaleEn represents en_US locale
	LocaleEn string = "en_US"
	// LocaleRu represents ru_RU locale
	LocaleRu string = "ru_RU"
)

// Translation is a representation of translation in Sylius
type Translation struct {
	Description string `json:"description,omitempty"`
	ID          int    `json:"id,omitempty"`
	Locale      string `json:"locale,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
}
