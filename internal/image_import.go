package internal

import (
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"
	"highlite-parser/internal/sylius/transfer"
	"context"
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
func (i *ImageImport) Import(product transfer.Product, h highlite.Product) error {
	product.Images = []transfer.Image{
		{
			Type: "test_upload",
		},
	}

	if err := i.client.TestImageUpload(context.Background(), product); err != nil {
		return err
	}

	return nil
}
