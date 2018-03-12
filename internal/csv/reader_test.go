package csv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	Separator    rune
	QuotedQuotes bool
	Input        string
	Output       [][]string
}{
	{
		Input:  "",
		Output: [][]string{},
	},
	{
		Input:  ",",
		Output: [][]string{{"", ""}},
	},
	{
		Input:  ",\n",
		Output: [][]string{{"", ""}},
	},
	{
		Input:  ",\n\n",
		Output: [][]string{{"", ""}, {""}},
	},
	{
		Input:  "\n\n\n",
		Output: [][]string{{""}, {""}, {""}},
	},
	{
		QuotedQuotes: true,
		Separator:    ';',
		Input: `field1;'field; ''2''';"field ""3"";
""field 3""
field 3";field4`,
		Output: [][]string{{"field1", "field; '2'", "field \"3\";\n\"field 3\"\nfield 3", "field4"}},
	},
	{
		Input:  `field1,'field, 世界, '2'',"field, "3""`,
		Output: [][]string{{"field1", "field, 世界, '2'", "field, \"3\""}},
	},
}

func TestParser(t *testing.T) {
	for _, tc := range testCases {
		// arrange
		parser := NewReader(strings.NewReader(tc.Input))
		parser.QuotedQuotes = tc.QuotedQuotes
		parser.FieldsFixed = false
		if tc.Separator != 0 {
			parser.Separator = tc.Separator
		}

		// act
		actual, err := readAll(parser)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, tc.Output, actual)
	}
}

func TestParserColumnsFixed(t *testing.T) {
	// arrange
	parser := NewReader(strings.NewReader("t1,t2\nt1,t2,t3\n"))

	// act
	actual, err := readAll(parser)

	// assert
	assert.EqualError(t, err, "wrong column count, row index 2")
	assert.Equal(t, [][]string{{"t1", "t2"}}, actual)
}

func TestParserColumnsFixedAssigned(t *testing.T) {
	// arrange
	parser := NewReader(strings.NewReader("t1,t2,t3\nt1,t2,t3\nt1,t2,t3,t4\nt1,t2,t3\n"))
	parser.FieldsCount = 2

	// act
	actual, err := readAll(parser)

	// assert
	assert.EqualError(t, err, "wrong column count, row index 1")
	assert.Equal(t, [][]string{}, actual)
}

func TestUnquotedQuote(t *testing.T) {
	// arrange
	parser := NewReader(strings.NewReader(`field1,field2,field3
field1,field "2","field "3" field 3"
field1,field2,field3
`))
	parser.QuotedQuotes = true

	// act
	actual, err := readAll(parser)

	// assert
	assert.EqualError(t, err, "unquoted quote on 2:25")
	assert.Equal(t, [][]string{{"field1", "field2", "field3"}}, actual)
}

func readAll(parser *Reader) ([][]string, error) {
	var actual = make([][]string, 0)
	for parser.Next() {
		actual = append(actual, parser.Values())
	}

	return actual, parser.Err()
}
