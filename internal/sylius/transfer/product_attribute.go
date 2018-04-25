package transfer

// IProductAttribute ...
type IProductAttribute interface {
	Equals(attr IProductAttribute) bool
}

// ProductAttributeRawHelper ...
type ProductAttributeRawHelper struct {
	Attribute  string      `json:"attribute"`
	Code       string      `json:"code"`
	LocaleCode string      `json:"localeCode"`
	Value      interface{} `json:"value"`
	Type       string      `json:"type"`
}

// ProductAttributeSelectSingle ...
type ProductAttributeSelectSingle struct {
	Code       string `json:"code,omitempty"`
	LocaleCode string `json:"localeCode"`
	Value      string `json:"value"`
}

// Equals ...
func (p ProductAttributeSelectSingle) Equals(attr IProductAttribute) bool {
	selectAttr, ok := attr.(ProductAttributeSelectSingle)
	if !ok {
		return false
	}

	if selectAttr.Code != p.Code {
		return false
	}

	if selectAttr.LocaleCode != p.LocaleCode {
		return false
	}

	if selectAttr.Value != p.Value {
		return false
	}

	return true
}
