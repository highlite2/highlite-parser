package imprt

import (
	"context"
	"strings"

	"highlite-parser/internal/cache"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/highlite/image"
	"highlite-parser/internal/highlite/translation"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"
	"highlite-parser/internal/sylius/transfer"
)

// NewProductImport creates new ProductImport.
func NewProductImport(client sylius.IClient, memo cache.IMemo,
	logger log.ILogger, dictionary translation.IDictionary, imageProvider image.IProvider) *ProductImport {
	return &ProductImport{
		logger:         logger,
		channelName:    "default", // TODO take it from config
		client:         client,
		categoryImport: NewCategoryImport(client, memo, logger),
		dictionary:     dictionary,
		imageProvider:  imageProvider,
	}
}

// ProductImport imports highlite product into sylius.
type ProductImport struct {
	logger         log.ILogger
	channelName    string
	client         sylius.IClient
	categoryImport *CategoryImport
	dictionary     translation.IDictionary
	imageProvider  image.IProvider
}

// Import imports highlite product into sylius.
func (i *ProductImport) Import(ctx context.Context, high highlite.Product) error {
	if product, err := i.client.GetProduct(ctx, high.Code); err == nil {
		return i.updateProduct(ctx, product, high)
	} else if err == sylius.ErrNotFound {
		return i.createProduct(ctx, high)
	} else {
		return err
	}
}

// Creates product.
func (i *ProductImport) createProduct(ctx context.Context, high highlite.Product) error {
	if _, err := i.categoryImport.Import(ctx, high.Category3); err != nil {
		return err
	}

	product := i.getProductFromHighlite(transfer.ProductEntire{}, high)

	images, imageErr := i.imageProvider.GetImages(ctx, high.Images)
	if imageErr != nil {
		return imageErr
	}

	defer func() {
		for _, reader := range images {
			reader.Close()
		}
	}()

	// TODO change image structure and distinguish the "main" image (cover image)
	productEntire, createErr := i.client.CreateProduct(ctx, product, images)
	if createErr != nil {
		return createErr
	}

	variant := i.getVariantFromHighlite(transfer.VariantEntire{}, high)
	if _, err := i.client.CreateProductVariant(ctx, productEntire.Code, variant); err != nil {
		return err
	}

	return nil
}

// Updates product.
// TODO make PATCH requests only if data really has changed!
func (i *ProductImport) updateProduct(ctx context.Context, product *transfer.ProductEntire, high highlite.Product) error {
	if err := i.client.UpdateProduct(ctx, i.getProductFromHighlite(*product, high)); err != nil {
		return err
	}

	variantCode := getProductMainVariantCode(product.Code)
	if variantEntire, err := i.client.GetProductVariant(ctx, product.Code, variantCode); err != nil {
		if err != sylius.ErrNotFound {
			return err
		}

		variant := i.getVariantFromHighlite(transfer.VariantEntire{}, high)
		if _, err := i.client.CreateProductVariant(ctx, product.Code, variant); err != nil {
			return err
		}
	} else {
		variant := i.getVariantFromHighlite(*variantEntire, high)
		if err := i.client.UpdateProductVariant(ctx, product.Code, variant); err != nil {
			return err
		}
	}

	return nil
}

// Creates sylius Variant structure from higlite product structure.
func (i *ProductImport) getVariantFromHighlite(variantEntire transfer.VariantEntire, high highlite.Product) transfer.Variant {
	variant := transfer.Variant{VariantEntire: variantEntire}

	variant.Code = getProductMainVariantCode(high.Code)
	variant.Translations = map[string]transfer.Translation{
		transfer.LocaleEn: {
			Name: high.Name,
		},
		transfer.LocaleRu: {
			Name: high.Name,
		},
	}
	variant.ChannelPrices = map[string]transfer.ChannelPrice{
		i.channelName: {
			Price: high.Price,
		},
	}
	variant.Width = high.Width
	variant.Height = high.Height
	variant.Weight = high.Weight
	variant.Depth = high.Length

	return variant
}

// Creates sylius Product structure from higlite product structure.
func (i *ProductImport) getProductFromHighlite(productEntire transfer.ProductEntire, high highlite.Product) transfer.Product {
	product := transfer.Product{ProductEntire: productEntire}
	product.Code = high.Code
	product.MainTaxon = high.Category3.GetCode()
	product.ProductTaxons = strings.Join([]string{high.Category3.GetCode(), high.Category2.GetCode(), high.Category1.GetCode()}, ",")
	product.Channels = []string{i.channelName}

	if len(product.Translations) == 0 {
		product.Translations = make(map[string]transfer.Translation)
	}

	product.Translations[transfer.LocaleEn] = transfer.Translation{
		Name:             high.Name,
		Slug:             high.URL,
		Description:      high.ProductDescription(),
		ShortDescription: high.SubHeading,
	}

	if item, ok := i.dictionary.Get(transfer.LocaleRu, high.No); ok {
		product.Translations[transfer.LocaleRu] = transfer.Translation{
			Name:             high.Name,
			Slug:             high.URL,
			Description:      item.GetDescription(),
			ShortDescription: item.GetShortDescription(),
		}
	} else {
		i.logger.Warnf("Can't find translations for product No %s", product.Code)
		product.Translations[transfer.LocaleRu] = transfer.Translation{
			Name:             high.Name,
			Slug:             high.URL,
			Description:      high.ProductDescription(),
			ShortDescription: high.SubHeading,
		}
	}

	product.Enabled = true

	switch high.Status {
	case highlite.StatusDecline, highlite.StatusEOL:
		product.Enabled = false
	}

	return product
}

func getProductMainVariantCode(productCode string) string {
	return productCode + "_main"
}
