package transfer

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
