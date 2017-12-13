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

// NewCategory creates new category and sets slug and code from category name.
func NewCategory(name string, parent *Category) *Category {
	cat := &Category{
		Name:   name,
		Parent: parent,
	}

	cat.SetSlugAndCode(name)

	return cat
}

// Category is a highlite import category.
type Category struct {
	Name   string
	Code   string
	Slug   string
	Parent *Category
}

// GetCode returns category full code (combined with parent's code).
func (c *Category) GetCode() string {
	if c.Parent != nil {
		return c.Parent.GetCode() + "_" + c.Code
	}

	return c.Code
}

// GetSlug returns category full slug (combined with parent's slug).
func (c *Category) GetSlug() string {
	if c.Parent != nil {
		return c.Parent.GetSlug() + "/" + c.Slug
	}

	return c.Slug
}

// SetSlugAndCode sets slug and code from given string.
func (c *Category) SetSlugAndCode(str string) {
	str = strings.ToLower(str)
	str = categoryCodeRegExp.ReplaceAllString(str, " ")
	fields := strings.Fields(str)

	c.Slug = strings.Join(fields, "-")
	c.Code = strings.Join(fields, "_")
}
