package imprt

import (
	"context"
	"fmt"
	"strings"

	"highlite2-import/internal/highlite"
	"highlite2-import/internal/highlite/image"
	"highlite2-import/internal/highlite/translation"
	"highlite2-import/internal/log"
	"highlite2-import/internal/sylius"
	"highlite2-import/internal/sylius/transfer"
)

// IProductImport imports highlite product into sylius.
type IProductImport interface {
	// Import imports highlite product into sylius.
	Import(ctx context.Context, high highlite.Product) error
}

// NewProductImport creates new ProductImport.
func NewProductImport(client sylius.IClient, categoryImport ICategoryImport,
	logger log.ILogger, dictionary translation.IDictionary, imageProvider image.IProvider) *ProductImport {
	return &ProductImport{
		logger:         logger,
		channelName:    "default", // TODO take it from config
		client:         client,
		categoryImport: categoryImport,
		dictionary:     dictionary,
		imageProvider:  imageProvider,
	}
}

// ProductImport imports highlite product into sylius.
type ProductImport struct {
	logger         log.ILogger
	channelName    string
	client         sylius.IClient
	categoryImport ICategoryImport
	dictionary     translation.IDictionary
	imageProvider  image.IProvider
}

// Import imports highlite product into sylius.
func (i *ProductImport) Import(ctx context.Context, high highlite.Product) error {
	if product, err := i.client.GetProduct(ctx, high.Code); err == nil {
		return i.updateProduct(ctx, *product, high)
	} else if err == sylius.ErrNotFound {
		return i.createProduct(ctx, high)
	} else {
		return fmt.Errorf("import: GetProduct client request returned error: %s", err)
	}
}

// Creates product.
func (i *ProductImport) createProduct(ctx context.Context, high highlite.Product) error {
	if _, err := i.categoryImport.Import(ctx, high.Category3); err != nil {
		return fmt.Errorf("createProduct: failed to import category: %s", err)
	}

	imageBucket, imageErr := i.imageProvider.GetImages(ctx, high.Images)
	if imageErr != nil {
		return fmt.Errorf("createProduct: failed to download images: %s", imageErr)
	}
	defer imageBucket.Close()

	productNew := i.getProductFromHighlite(high)
	images := prepareImages(high, &productNew, imageBucket)

	productEntire, createErr := i.client.CreateProduct(ctx, productNew, images)
	if createErr != nil {
		return fmt.Errorf("createProduct: client CreateProduct returned error: %s", createErr)
	}

	variantNew := i.getVariantFromHighlite(high)
	if _, err := i.client.CreateProductVariant(ctx, productEntire.Code, variantNew); err != nil {
		return fmt.Errorf("createProduct: client CreateProductVariant returned error: %s", err)
	}

	return nil
}

// Updates product.
func (i *ProductImport) updateProduct(ctx context.Context, productEntire transfer.ProductEntire, high highlite.Product) error {
	productNew := i.getProductFromHighlite(high)
	if !transfer.ProductsEqual(productEntire, productNew) {
		i.logger.Infof("Updating product %s", productEntire.Code)
		if err := i.client.UpdateProduct(ctx, productNew); err != nil {
			return fmt.Errorf("updateProduct: client UpdateProduct returned error: %s", err)
		}
	}

	variantNew := i.getVariantFromHighlite(high)
	if variantEntire, err := i.client.GetProductVariant(ctx, productEntire.Code, getProductMainVariantCode(productEntire.Code)); err != nil {
		if err != sylius.ErrNotFound {
			return fmt.Errorf("updateProduct: client GetProductVariant returned error: %s", err)
		}

		if _, err := i.client.CreateProductVariant(ctx, productEntire.Code, variantNew); err != nil {
			return fmt.Errorf("updateProduct: client CreateProductVariant returned error: %s", err)
		}
	} else if !transfer.VariantsEqual(*variantEntire, variantNew) {
		i.logger.Infof("Updating variant %s", getProductMainVariantCode(productEntire.Code))
		if err := i.client.UpdateProductVariant(ctx, productEntire.Code, variantNew); err != nil {
			return fmt.Errorf("updateProduct: client UpdateProductVariant returned error: %s", err)
		}
	}

	return nil
}

// Creates sylius Variant structure from higlite product structure.
func (i *ProductImport) getVariantFromHighlite(high highlite.Product) transfer.Variant {
	return transfer.Variant{
		VariantEntire: transfer.VariantEntire{
			Code: getProductMainVariantCode(high.Code),
			Translations: map[string]transfer.Translation{
				transfer.LocaleEn: {
					Name: strings.Trim(high.Name, "\n "),
				},
				transfer.LocaleRu: {
					Name: strings.Trim(high.Name, "\n "),
				},
			},
			ChannelPrices: map[string]transfer.ChannelPrice{
				i.channelName: {
					Price: high.Price,
				},
			},
		},
	}
}

// Creates sylius Product structure from higlite product structure.
func (i *ProductImport) getProductFromHighlite(high highlite.Product) transfer.Product {
	tr := transfer.Translation{
		Name:             high.ProductName(),
		Slug:             high.URL,
		Description:      high.GetDescription(),
		ShortDescription: high.GetShortDescription(),
	}

	product := transfer.Product{
		ProductEntire: transfer.ProductEntire{
			Code:    high.Code,
			Enabled: true,
			Translations: map[string]transfer.Translation{
				transfer.LocaleEn: tr,
				transfer.LocaleRu: tr,
			},
		},
		MainTaxon: high.Category3.GetCode(),
		ProductTaxons: strings.Join(
			[]string{
				high.Category3.GetCode(),
				high.Category2.GetCode(),
				high.Category1.GetCode(),
				high.CategoryRoot.GetCode(),
			},
			",",
		),
		Channels: []string{i.channelName},
	}

	if item, ok := i.dictionary.Get(transfer.LocaleRu, high.No); ok && !item.Empty() {
		tr.Description = item.GetDescription()
		tr.ShortDescription = item.GetShortDescription()
		product.Translations[transfer.LocaleRu] = tr
	} else {
		i.logger.Warnf("Can't find translations for product No %s", product.Code)
	}

	return product
}

// Product main variant name.
func getProductMainVariantCode(productCode string) string {
	return productCode + "_main"
}

// Prepares image structure to pass it to sylius API.
func prepareImages(high highlite.Product, product *transfer.Product, bucket image.Bucket) []transfer.ImageUpload {
	var images []transfer.ImageUpload

	if len(bucket) > 0 {
		product.Images = make([]transfer.Image, len(bucket))
		images = make([]transfer.ImageUpload, len(bucket))

		for i, item := range bucket {
			if item.Name == high.Images[0] {
				product.Images[i] = transfer.Image{Type: "main"}
			}

			images[i] = transfer.ImageUpload{
				Name:   item.Name,
				Reader: item.Reader,
			}
		}
	}

	return images
}
