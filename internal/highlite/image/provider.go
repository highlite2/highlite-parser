package image

import (
	"context"
	"io"
	"net/http"
)

// IProvider is an interface, that is supposed to return highlite product images.
type IProvider interface {
	GetImages(ctx context.Context, images []string) (map[string]io.ReadCloser, error)
}

var _ IProvider = (*HTTPProvider)(nil)

// HTTPProvider is a implementation of IProvider.
type HTTPProvider struct {
	imageGet func(string) (*http.Response, error)
}

// GetImages downloads images from the internet.
func (h HTTPProvider) GetImages(ctx context.Context, images []string) (map[string]io.ReadCloser, error) {
	internal := &httpReader{
		downloadFn:   http.Get,
		downloads:    make(chan downloadResponse),
		ready:        make(chan bool),
		imageReaders: make(map[string]io.ReadCloser),
		imageNames:   images,
	}

	return internal.downloadImages(ctx)
}
