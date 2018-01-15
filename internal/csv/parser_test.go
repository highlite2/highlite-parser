package csv

import (
	"strings"
	"testing"

	"io"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	Input  string
	Output [][]string
}{
	{
		Input:  "",
		Output: [][]string{},
	},
	{
		Input: `,
`,
		Output: [][]string{
			{"", ""},
		},
	},
	{
		Input: `,`,
		Output: [][]string{
			{"", ""},
		},
	},
	{
		Input: "t1,t2,t3",
		Output: [][]string{
			{"t1", "t2", "t3"},
		},
	},
	{
		Input: `t1,t2,t3`,
		Output: [][]string{
			{"t1", "t2", "t3"},
		},
	},
	{
		Input: `,t1,t2,t3,,,t4,
`,
		Output: [][]string{
			{"", "t1", "t2", "t3", "", "", "t4", ""},
		},
	},
	{
		Input: `"pa""rs"er",'parser',parser`,
		Output: [][]string{
			{"pa\"rs\"er", "parser", "parser"},
		},
	},
	{
		Input: `"pa
rs
er","pa"rser",pa"rs"er`,
		Output: [][]string{
			{"pa\nrs\ner", "pa\"rser", "pa\"rs\"er"},
		},
	},
}

func TestParser(t *testing.T) {
	for _, tc := range testCases {
		// arrange
		parser := NewParser(strings.NewReader(tc.Input))

		// act
		var actual [][]string = make([][]string, 0)
		var titles []string
		var err error
		for {
			titles, err = parser.Next()
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}

			actual = append(actual, titles)
		}

		// assert
		assert.Nil(t, err)
		assert.Equal(t, tc.Output, actual)
	}

}
