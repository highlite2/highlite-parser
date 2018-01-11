package highlite

import "fmt"

// IDictionary is a service to get translations for different languages.
type IDictionary interface {
	Get(lang string, id string) (ITranslation, bool)
}

// ITranslation is a product translation for exact language
type ITranslation interface {
	ProductDescription() string
}

// DictionaryMap is a service to get translations for different languages.
type DictionaryMap struct {
	translations map[string]map[string]Translation
}

// Get returns a product translation for exact language.
func (t *DictionaryMap) Get(lang string, id string) (ITranslation, bool) {
	dictionary, ok := t.translations[lang]
	if !ok {
		return nil, false
	}

	translation, ok := dictionary[id]
	if !ok {
		return nil, false
	}

	return &translation, ok
}

// Translation ... TODO
type Translation struct {
	CatTextMain string
	CatTextSubH string
	USP         string
	TechSpec    string
}

// ProductDescription return product description from several fields
func (t *Translation) ProductDescription() string {
	description := ""
	description += replaceHTMLEntities(t.USP)
	description += "\n"
	description += replaceHTMLEntities(t.CatTextMain)
	description += "\n\n"
	description += replaceHTMLEntities(t.TechSpec)

	return description
}
