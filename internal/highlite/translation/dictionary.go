package translation

// IDictionary is a service to get translations for different languages.
type IDictionary interface {
	Get(lang string, id string) (IItem, bool)
}

// IItem is a product translation for exact language.
type IItem interface {
	GetDescription() string
}

// NewMemoryDictionary creates a new instance of a MemoryDictionary.
func NewMemoryDictionary() *MemoryDictionary {
	return &MemoryDictionary{
		languages: make(map[string]map[string]IItem),
	}
}

// MemoryDictionary is an implementation of IDictionary.
// It contains all translations in memory.
type MemoryDictionary struct {
	languages map[string]map[string]IItem
}

// Get returns a product translation for exact language.
func (t *MemoryDictionary) Get(lang string, id string) (IItem, bool) {
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
