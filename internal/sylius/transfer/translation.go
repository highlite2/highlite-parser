package transfer

const (
	// LocaleEn represents en_US locale
	LocaleEn string = "en_US"
	// LocaleRu represents ru_RU locale
	LocaleRu string = "ru_RU"
)

// Translation is a representation of translation in Sylius
type Translation struct {
	Name             string `json:"name,omitempty"`
	Slug             string `json:"slug,omitempty"`
	Description      string `json:"description,omitempty"`
	ShortDescription string `json:"shortDescription,omitempty"`
}
