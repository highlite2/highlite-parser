package highlite

import (
	"regexp"
	"strings"
)

// Regexp to cut all symbols except a-z0-9.
var categoryCodeRegExp *regexp.Regexp

// Initialisation for the regexp. Panics if failed.
func init() {
	reg, err := regexp.Compile("[^a-z0-9]+")
	if err != nil {
		panic(err)
	}

	categoryCodeRegExp = reg
}

// Identification is a helper for setting code and url of the resource.
type Identification struct {
	Code string
	URL  string
}

// SetCodeAndURL sets url and code from given string.
func (i *Identification) SetCodeAndURL(str string) {
	str = strings.ToLower(str)
	str = categoryCodeRegExp.ReplaceAllString(str, " ")
	fields := strings.Fields(str)

	i.URL = strings.Join(fields, "-")
	i.Code = strings.Join(fields, "_")
}
