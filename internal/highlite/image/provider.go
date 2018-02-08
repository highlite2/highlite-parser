package image

import (
	"context"
	"io"
	"net/http"
)

const highliteImageLocation = "http://www.highlite.nl/var/StorageHighlite/ProduktBilder/"

// Bucket is an image readers map.
type Bucket map[string]io.ReadCloser

// Close image readers.
func (b Bucket) Close() {
	for _, reader := range b {
		reader.Close()
	}
}

// IProvider is an interface, that is supposed to return highlite product images.
type IProvider interface {
	GetImages(ctx context.Context, images []string) (Bucket, error)
}

var _ IProvider = (*HTTPProvider)(nil)

// HTTPProvider is a implementation of IProvider.
type HTTPProvider struct {
	imageGet func(string) (*http.Response, error)
}

// GetImages loadsChan images from the internet.
func (h HTTPProvider) GetImages(ctx context.Context, images []string) (Bucket, error) {
	internal := &httpReader{}
	internal.init(http.Get, images, highliteImageLocation)

	return internal.downloadImages(ctx)
}
