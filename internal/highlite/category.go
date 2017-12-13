package highlite

import (
	"regexp"
	"strings"
)

var categoryCodeRegExp *regexp.Regexp

func init() {
	reg, err := regexp.Compile("[^a-z0-9]+")
	if err != nil {
		panic(err)
	}

	categoryCodeRegExp = reg
}

// NewCategory creates new category and sets url and code from category name.
func NewCategory(name string, parent *Category) *Category {
	cat := &Category{Name: name, Parent: parent}
	cat.SetCodeAndURL(name)

	return cat
}

// Category is a highlite import category.
type Category struct {
	Name   string
	Code   string
	URL    string
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
	if c.Parent != nil {
		return c.Parent.GetURL() + "/" + c.URL
	}

	return c.URL
}

// SetCodeAndURL sets url and code from given string.
func (c *Category) SetCodeAndURL(str string) {
	str = strings.ToLower(str)
	str = categoryCodeRegExp.ReplaceAllString(str, " ")
	fields := strings.Fields(str)

	c.URL = strings.Join(fields, "-")
	c.Code = strings.Join(fields, "_")
}
