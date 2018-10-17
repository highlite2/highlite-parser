package translation

// IDictionary is a service to get translations for different languages.
type IDictionary interface {
	Get(lang string, id string) (IProduct, bool)
}

// IProduct is a product translation for exact language.
type IProduct interface {
	GetDescription() string
	GetShortDescription() string
}

// NewMemoryDictionary creates a new instance of a MemoryDictionary.
func NewMemoryDictionary() *MemoryDictionary {
	return &MemoryDictionary{
		languages: make(map[string]map[string]IProduct),
	}
}

// MemoryDictionary is an implementation of IDictionary.
// It contains all translations in memory.
type MemoryDictionary struct {
	languages map[string]map[string]IProduct
}

// Get returns a product translation for exact language.
func (t *MemoryDictionary) Get(lang string, id string) (IProduct, bool) {
	language, ok := t.languages[lang]
	if !ok {
		return nil, false
	}

	product, ok := language[id]
	if !ok {
		return nil, false
	}

	return product, ok
}

// GetMap returns dictionary map
func (t *MemoryDictionary) GetMap() map[string]map[string]IProduct {
	return t.languages
}
