package csv

import (
	"fmt"
	"strconv"
)

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

// CheckTitles checks if all titles exist in a title map
func (t *TitleMap) CheckTitles(titles []string) error {
	for _, title := range titles {
		if _, ok := t.titles[title]; !ok {
			return fmt.Errorf("can not find %s in title map", title)
		}
	}

	return nil
}

// CheckValues checks if values count equals titles count
func (t *TitleMap) CheckValues(values []string) error {
	if len(t.titles) != len(values) {
		return fmt.Errorf("wrong values count: %d, must be: %d", len(values), len(t.titles))
	}

	return nil
}

// GetString finds a value from values slice by title. First, it gets title index
// from title map, and then gets the string value from values slice by that index.
func (t *TitleMap) GetString(title string, values []string) string {
	i, ok := t.titles[title]
	if ok && i < len(values) {
		return t.processValue(values[i])
	}

	return ""
}

// GetFloat finds a value from values slice by title. First, it gets title index
// from title map, and then gets the float64 value from values slice by that index.
func (t *TitleMap) GetFloat(title string, values []string) float64 {
	str := t.GetString(title, values)
	val, _ := strconv.ParseFloat(str, 64)

	return val
}

// GetInt finds a value from values slice by title. First, it gets title index
// from title map, and then gets the int64 value from values slice by that index.
func (t *TitleMap) GetInt(title string, values []string) int64 {
	str := t.GetString(title, values)
	val, _ := strconv.ParseInt(str, 10, 64)

	return val
}

// Calls a callback if it is set.
func (t *TitleMap) processValue(v string) string {
	if t.callback != nil {
		return t.callback(v)
	}

	return v
}
