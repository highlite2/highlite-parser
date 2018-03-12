package image

import (
	"context"
	"io"
	"net/http"
)

const highliteImageLocation = "http://www.highlite.nl/var/StorageHighlite/ProduktBilder/"

// Bucket is a collection of readers.
type Bucket []BucketItem

// BucketItem is an image reader.
type BucketItem struct {
	Name   string
	Reader io.Reader
}

// Close closes image readers reader implements io.Close interface.
func (b Bucket) Close() {
	for _, item := range b {
		if closer, ok := item.Reader.(io.Closer); ok {
			closer.Close()
		}
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
