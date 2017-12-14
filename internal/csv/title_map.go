package csv

// NewTitleMap creates new TitleMap object and executes
// titles initialization
func NewTitleMap(titles []string) *TitleMap {
	m := &TitleMap{}
	m.Init(titles)

	return m
}

// TitleMap helps to get a value from csv row by column title.
type TitleMap struct {
	titles   map[string]int
	callback func(string) string
}

// SetCallback sets a callback, that will be invoked for every value
func (t *TitleMap) SetCallback(p func(string) string) {
	t.callback = p
}

// Init creates a title map: title => index
func (t *TitleMap) Init(titles []string) {
	t.titles = make(map[string]int)
	for i, v := range titles {
		t.titles[t.processValue(v)] = i
	}
}

// Get finds a value from values slice by title. First, it gets title index
// from title map, and then gets the value from values slice by that index.
func (t *TitleMap) Get(title string, values []string) string {
	i, ok := t.titles[title]
	if ok && i < len(values) {
		return t.processValue(values[i])
	}

	return ""
}

// Calls a callback if it is set.
func (t *TitleMap) processValue(v string) string {
	if t.callback != nil {
		return t.callback(v)
	}

	return v
}
