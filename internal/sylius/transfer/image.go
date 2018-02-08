package transfer

import "io"

// Image is a representation of a product image
type Image struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

// UploadImage is a help struct to upload images.
type ImageUpload struct {
	Name   string
	Reader io.Reader
}
