package transfer

// Image is a representation of a product image
type Image struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}
