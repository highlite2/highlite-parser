package transfer

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttribute(t *testing.T) {
	actual := getActualAttributeFromFile()
	expected := getAttribute()

	assert.Equal(t, expected, actual)
}

func getAttribute() Attribute {
	return Attribute{
		Code: "highlite_brand",
		Configuration: AttributeConfiguration{
			Choices: map[string]AttributeConfigurationChoice{
				"highlite_brand_option_1": {
					"en_US": "Brand 1",
					"ru_RU": "Бренд 1",
				},
				"highlite_brand_option_2": {
					"en_US": "Brand 2",
					"ru_RU": "Бренд 2",
				},
				"highlite_brand_option_3": {
					"en_US": "Brand 3",
					"ru_RU": "Бренд 3",
				},
			},
		},
		Translations: map[string]Translation{
			"en_US": {
				Name: "Brand",
			},
			"ru_RU": {
				Name: "Бренд",
			},
		},
	}
}

func getActualAttributeFromFile() Attribute {
	bytes, err := ioutil.ReadFile("_test_data/attribute-type-select-example.json")
	if err != nil {
		panic(err)
	}

	attribute := Attribute{}
	err = json.Unmarshal(bytes, &attribute)
	if err != nil {
		panic(err)
	}

	return attribute
}
