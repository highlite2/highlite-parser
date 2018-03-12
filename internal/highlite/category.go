package highlite

// NewCategory creates new category and sets url and code from category name.
func NewCategory(name string, parent *Category) *Category {
	cat := &Category{Name: name, Parent: parent}
	cat.SetCodeAndURL(name)

	return cat
}

// Category is a highlite import category.
type Category struct {
	Identification

	Root   bool
	Name   string
	Parent *Category
}

// GetCode returns category full code (combined with parent's code).
func (c *Category) GetCode() string {
	if c.Parent != nil {
		return c.Parent.GetCode() + "_" + c.Code
	}

	return c.Code
}

// GetURL returns category full url (combined with parent's url).
func (c *Category) GetURL() string {
	if c.Parent != nil && !c.Parent.Root {
		return c.Parent.GetURL() + "/" + c.URL
	}

	return c.URL
}
