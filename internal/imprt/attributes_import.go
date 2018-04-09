package imprt

import (
	"context"
	"fmt"
	"sync"

	"highlite2-import/internal/highlite"
	"highlite2-import/internal/log"
	"highlite2-import/internal/sylius"
	"highlite2-import/internal/sylius/transfer"
)

const attributeBrandCode = "highlite-brand"

// IAttributesImport is a Sylius attributes importer.
type IAttributesImport interface {
	SetProductAttributes(ctx context.Context, high highlite.Product, product *transfer.Product) error
}

// NewAttributesImport creates new Sylius attributes importer.
func NewAttributesImport(client sylius.IClient, logger log.ILogger) *AttributesImport {
	return &AttributesImport{
		logger:        logger,
		client:        client,
		attrBrandLock: sync.RWMutex{},
	}
}

// AttributesImport is a Sylius attributes importer implementation.
type AttributesImport struct {
	logger        log.ILogger
	client        sylius.IClient
	attrBrand     *transfer.Attribute
	attrBrandLock sync.RWMutex
}

// Init initializes actual for import attributes.
func (i *AttributesImport) Init(ctx context.Context) error {
	return i.initBrandAttribute(ctx)
}

// SetProductAttributes sets product attributes.
func (i *AttributesImport) SetProductAttributes(ctx context.Context, high highlite.Product, product *transfer.Product) error {
	brand, err := i.getBrandAttributeChoiceCode(ctx, high)
	if err != nil {
		return err
	}

	product.Attributes = []transfer.ProductAttribute{
		{
			Attribute:  attributeBrandCode,
			LocaleCode: transfer.LocaleEn,
			Value:      brand,
		},
		{
			Attribute:  attributeBrandCode,
			LocaleCode: transfer.LocaleRu,
			Value:      brand,
		},
	}

	return nil
}

// Initializes brand attribute.
func (i *AttributesImport) initBrandAttribute(ctx context.Context) error {
	var err error

	i.attrBrand, err = i.client.GetProductAttribute(ctx, attributeBrandCode)
	if err != nil {
		if err != sylius.ErrNotFound {
			return fmt.Errorf("fetching brand attribute error: %s", err)
		}

		i.attrBrand, err = i.client.CreateProductAttribute(ctx, transfer.AttributeTypeSelect, transfer.Attribute{
			Code: attributeBrandCode,
			Translations: map[string]transfer.Translation{
				transfer.LocaleRu: {Name: "Бренд"}, // TODO these translations must be moved to a common translation place.
				transfer.LocaleEn: {Name: "Brand"},
			},
		})

		if err != nil {
			return fmt.Errorf("creating brand attribute error: %s", err)
		}
	}

	return nil
}

// Gets brand attribute choice code by highlite product info.
func (i *AttributesImport) getBrandAttributeChoiceCode(ctx context.Context, high highlite.Product) (string, error) {
	if i.checkAttributeBrandChoiceCodeExists(high) {
		return high.GetBrandCode(), nil
	}

	if err := i.addAttributeBrandChoiceCode(ctx, high); err != nil {
		return "", fmt.Errorf("creating brand attribute choice error: %s", err)
	}

	return high.GetBrandCode(), nil
}

// Checks if brand attribute choice already exists.
func (i *AttributesImport) checkAttributeBrandChoiceCodeExists(high highlite.Product) bool {
	i.attrBrandLock.RLock()
	defer i.attrBrandLock.RUnlock()

	_, ok := i.attrBrand.Configuration.Choices[high.GetBrandCode()]

	return ok
}

// Adds new brand attribute choice to the brand attribute.
func (i *AttributesImport) addAttributeBrandChoiceCode(ctx context.Context, high highlite.Product) error {
	i.attrBrandLock.Lock()
	defer i.attrBrandLock.Unlock()

	if _, ok := i.attrBrand.Configuration.Choices[high.GetBrandCode()]; ok {
		return nil
	}

	if len(i.attrBrand.Configuration.Choices) == 0 {
		i.attrBrand.Configuration.Choices = make(map[string]transfer.AttributeConfigurationChoice)
	}

	i.attrBrand.Configuration.Choices[high.GetBrandCode()] = transfer.AttributeConfigurationChoice{
		transfer.LocaleRu: high.Brand,
		transfer.LocaleEn: high.Brand,
	}

	err := i.client.UpdateProductAttribute(ctx, *i.attrBrand)
	if err != nil {
		delete(i.attrBrand.Configuration.Choices, high.GetBrandCode())

		return err
	}

	return nil
}
