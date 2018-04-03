package transfer

import "encoding/json"

// AttributeType is an attribute type.
type AttributeType string

const (
	// AttributeTypeSelect Is a "Select" type of a product attribute.
	AttributeTypeSelect AttributeType = "select"
)

// Attribute is a representation of Sylius attribute.
type Attribute struct {
	Code          string                 `json:"code,omitempty"`
	Translations  map[string]Translation `json:"translations,omitempty"`
	Configuration AttributeConfiguration `json:"configuration,omitempty"`
}

// AttributeConfiguration is an attribute configuration.
type AttributeConfiguration struct {
	Choices map[string]AttributeConfigurationChoice `json:"choices,omitempty"`
}

// AttributeConfigurationChoice is a select attribute choice.
type AttributeConfigurationChoice map[string]string

// This struct helps to unmarshal json. Sylius api returns different
// representation for Choices: if there are no choices, Sylius returns
// an empty array, if not empty - an object.
type attributeConfigurationRaw struct {
	Choices json.RawMessage `json:"choices"`
}

// UnmarshalJSON helps to fix inconsistency in sylius api response.
func (t *AttributeConfiguration) UnmarshalJSON(value []byte) error {
	raw := &attributeConfigurationRaw{}
	if err := json.Unmarshal(value, raw); err != nil {
		return err
	}

	var ch map[string]AttributeConfigurationChoice
	if err := json.Unmarshal(raw.Choices, &ch); err == nil {
		t.Choices = ch
	}

	return nil
}
