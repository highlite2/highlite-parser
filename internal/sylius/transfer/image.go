package transfer

// Image is a representation of a product image
type Image struct {
	ID   int    `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}
