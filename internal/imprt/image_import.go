package imprt

import (
	"context"

	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"
	"highlite-parser/internal/sylius/transfer"
)

// NewImageImport creates new ImageImport.
func NewImageImport(client sylius.IClient, logger log.ILogger) *ImageImport {
	return &ImageImport{
		client: client,
		logger: logger,
	}
}

// ImageImport imports product images into sylius.
type ImageImport struct {
	client sylius.IClient
	logger log.ILogger
}

// Import TODO
func (i *ImageImport) Import(ctx context.Context, product transfer.Product, h highlite.Product) error {
	return nil
}
