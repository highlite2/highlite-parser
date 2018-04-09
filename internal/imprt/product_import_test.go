package imprt

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"highlite2-import/internal/highlite"
	img "highlite2-import/internal/highlite/image"
	"highlite2-import/internal/highlite/translation"
	"highlite2-import/internal/sylius"
	"highlite2-import/internal/sylius/transfer"
	"highlite2-import/internal/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct(t *testing.T) {
	// arrange
	high := getHighliteProductMock()
	highTranslate := getHighliteTranslationMock(high)

	client := &mocks.SyliusClient{}
	client.On("GetProduct", mock.Anything, high.Code).Once().Return(nil, sylius.ErrNotFound)
	client.On("CreateProduct", mock.Anything, getProductFromHighliteMock(high), getImageUploadMock(high.Images)).Once().Return(&transfer.ProductEntire{Code: high.Code}, nil)
	client.On("CreateProductVariant", mock.Anything, high.Code, getProductVariantFromHighliteMock(high)).Once().Return(nil, nil)

	categoryImport := &mocks.CategoryImport{}
	categoryImport.On("Import", mock.Anything, mock.Anything).Once().Return(nil, nil)

	logger := &mocks.Logger{}
	logger.On("Infof")

	dictionary := &mocks.TranslationDictionay{}
	dictionary.On("Get", transfer.LocaleRu, high.No).Once().Return(&highTranslate, true)

	imgProvider := &mocks.ImageProvider{}
	imgProvider.On("GetImages", mock.Anything, high.Images).Once().Return(getImageBucketMock(high.Images), nil)

	attrImport := &mocks.AttributesImport{}
	attrImport.On("SetProductAttributes", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)

	// act
	im := NewProductImport(client, categoryImport, logger, dictionary, imgProvider, attrImport)
	err := im.Import(context.Background(), high)

	// assert
	assert.Nil(t, err)
}

func TestUpdateProduct_NoUpdates(t *testing.T) {
	// arrange
	high := getHighliteProductMock()
	highTranslate := getHighliteTranslationMock(high)
	variant := getProductVariantFromHighliteMock(high)

	client := &mocks.SyliusClient{}
	client.On("GetProduct", mock.Anything, high.Code).Once().Return(getProductEntireFromHighliteMock(high), nil)
	client.On("GetProductVariant", mock.Anything, high.Code, variant.Code).Once().Return(getProductVariantEntireFromHighliteMock(high), nil)

	categoryImport := &mocks.CategoryImport{}

	logger := &mocks.Logger{}
	logger.On("Infof")

	dictionary := &mocks.TranslationDictionay{}
	dictionary.On("Get", transfer.LocaleRu, high.No).Once().Return(&highTranslate, true)

	imgProvider := &mocks.ImageProvider{}

	attrImport := &mocks.AttributesImport{}

	// act
	im := NewProductImport(client, categoryImport, logger, dictionary, imgProvider, attrImport)
	err := im.Import(context.Background(), high)

	// assert
	assert.Nil(t, err)
}

func TestUpdateProduct_ProductUpdate(t *testing.T) {
	// arrange
	high := getHighliteProductMock()
	highTranslate := getHighliteTranslationMock(high)
	variant := getProductVariantFromHighliteMock(high)
	product := getProductFromHighliteMock(high)
	product.Images = nil

	productEntire := getProductEntireFromHighliteMock(high)
	tr := productEntire.Translations[transfer.LocaleRu]
	tr.Name = tr.Name + " Sale!"
	productEntire.Translations[transfer.LocaleRu] = tr

	client := &mocks.SyliusClient{}
	client.On("GetProduct", mock.Anything, high.Code).Once().Return(productEntire, nil)
	client.On("UpdateProduct", mock.Anything, product).Once().Return(nil)
	client.On("GetProductVariant", mock.Anything, high.Code, variant.Code).Once().Return(getProductVariantEntireFromHighliteMock(high), nil)

	categoryImport := &mocks.CategoryImport{}
	logger := &mocks.Logger{}
	logger.On("Infof", mock.Anything, mock.Anything).Maybe()

	dictionary := &mocks.TranslationDictionay{}
	dictionary.On("Get", transfer.LocaleRu, high.No).Once().Return(&highTranslate, true)
	imgProvider := &mocks.ImageProvider{}

	attrImport := &mocks.AttributesImport{}

	// act
	im := NewProductImport(client, categoryImport, logger, dictionary, imgProvider, attrImport)
	err := im.Import(context.Background(), high)

	// assert
	assert.Nil(t, err)
}

func TestUpdateProduct_VariantUpdate(t *testing.T) {
	// arrange
	high := getHighliteProductMock()
	highTranslate := getHighliteTranslationMock(high)
	variant := getProductVariantFromHighliteMock(high)
	product := getProductFromHighliteMock(high)
	product.Images = nil

	variantEntire := getProductVariantEntireFromHighliteMock(high)
	tr := variantEntire.Translations[transfer.LocaleRu]
	tr.Name = tr.Name + " Sale!"
	variantEntire.Translations[transfer.LocaleRu] = tr

	client := &mocks.SyliusClient{}
	client.On("GetProduct", mock.Anything, high.Code).Once().Return(getProductEntireFromHighliteMock(high), nil)
	client.On("GetProductVariant", mock.Anything, high.Code, variant.Code).Once().Return(variantEntire, nil)
	client.On("UpdateProductVariant", mock.Anything, high.Code, getProductVariantFromHighliteMock(high)).Once().Return(nil)

	categoryImport := &mocks.CategoryImport{}
	logger := &mocks.Logger{}
	logger.On("Infof", mock.Anything, mock.Anything).Maybe()

	dictionary := &mocks.TranslationDictionay{}
	dictionary.On("Get", transfer.LocaleRu, high.No).Once().Return(&highTranslate, true)
	imgProvider := &mocks.ImageProvider{}

	attrImport := &mocks.AttributesImport{}

	// act
	im := NewProductImport(client, categoryImport, logger, dictionary, imgProvider, attrImport)
	err := im.Import(context.Background(), high)

	// assert
	assert.Nil(t, err)
}

func getImageBucketMock(images []string) img.Bucket {
	var bucket = make(img.Bucket, len(images))
	for i, im := range images {
		bucket[i] = img.BucketItem{
			Name:   im,
			Reader: bytes.NewReader([]byte(im)),
		}
	}

	return bucket
}

func getImageUploadMock(images []string) []transfer.ImageUpload {
	uploads := make([]transfer.ImageUpload, len(images))
	for i, im := range images {
		uploads[i] = transfer.ImageUpload{
			Name:   im,
			Reader: bytes.NewReader([]byte(im)),
		}
	}

	return uploads
}

func getHighliteProductMock() highlite.Product {
	cat0 := highlite.NewCategory("Root", nil)
	cat0.Root = true
	cat1 := highlite.NewCategory("Category1", cat0)
	cat2 := highlite.NewCategory("Category2", cat1)
	cat3 := highlite.NewCategory("Category3", cat2)

	product := highlite.Product{
		No:          "No",
		Name:        "Name",
		SubHeading:  "SubHeading",
		Description: "Description",
		Specs:       "Specs",
		Brand:       "Brand",

		InStock: "InStock",
		Status:  "Status",

		Country: "Country",
		Price:   123.45,

		CategoryRoot: cat0,
		Category1:    cat1,
		Category2:    cat2,
		Category3:    cat3,

		Images: []string{"image1.jpg", "image2.jpg"},
	}

	product.SetCodeAndURL(product.Name + " " + product.No)
	product.Code = "highlite-" + product.No

	return product
}

func getHighliteTranslationMock(high highlite.Product) translation.ProductCSV {
	return translation.ProductCSV{
		MainText:   high.Description,
		SubHeading: high.SubHeading,
		TechSpec:   high.Specs,
	}
}

func getProductFromHighliteMock(high highlite.Product) transfer.Product {
	highTranslation := getHighliteTranslationMock(high)

	product := transfer.Product{
		ProductEntire: transfer.ProductEntire{
			Code:    high.Code,
			Enabled: true,
			Translations: map[string]transfer.Translation{
				transfer.LocaleEn: {
					Name:             high.GetProductName(),
					Slug:             high.URL,
					Description:      high.GetProductDescription(),
					ShortDescription: high.SubHeading,
				},
				transfer.LocaleRu: {
					Name:             high.GetProductName(),
					Slug:             high.URL,
					Description:      highTranslation.GetDescription(),
					ShortDescription: highTranslation.GetShortDescription(),
				},
			},
		},
		MainTaxon:     high.Category3.GetCode(),
		ProductTaxons: strings.Join([]string{high.Category3.GetCode(), high.Category2.GetCode(), high.Category1.GetCode(), high.CategoryRoot.GetCode()}, ","),
		Channels:      []string{"default"},
	}

	if len(high.Images) > 0 {
		product.Images = []transfer.Image{{Type: "main"}, {}}
	}

	return product
}

func getProductEntireFromHighliteMock(high highlite.Product) *transfer.ProductEntire {
	highTranslation := getHighliteTranslationMock(high)

	return &transfer.ProductEntire{
		Code:    high.Code,
		Enabled: true,
		Translations: map[string]transfer.Translation{
			transfer.LocaleEn: {
				Name:             high.GetProductName(),
				Slug:             high.URL,
				Description:      high.GetProductDescription(),
				ShortDescription: high.SubHeading,
			},
			transfer.LocaleRu: {
				Name:             high.GetProductName(),
				Slug:             high.URL,
				Description:      highTranslation.GetDescription(),
				ShortDescription: highTranslation.GetShortDescription(),
			},
		},
		Images:   []transfer.Image{{Type: "main"}, {}},
		Channels: []transfer.Channel{{Code: "default"}},
		MainTaxon: &transfer.Taxon{
			Code: high.Category3.GetCode(),
		},
		ProductTaxons: []transfer.TaxonWrap{
			{Taxon: transfer.Taxon{Code: high.Category3.GetCode()}},
			{Taxon: transfer.Taxon{Code: high.Category2.GetCode()}},
			{Taxon: transfer.Taxon{Code: high.Category1.GetCode()}},
			{Taxon: transfer.Taxon{Code: high.CategoryRoot.GetCode()}},
		},
	}

}

func getProductVariantFromHighliteMock(high highlite.Product) transfer.Variant {
	return transfer.Variant{
		VariantEntire: transfer.VariantEntire{
			Code: high.Code + "_main",
			Translations: map[string]transfer.Translation{
				transfer.LocaleEn: {
					Name: high.Brand + " " + high.Name,
				},
				transfer.LocaleRu: {
					Name: high.Brand + " " + high.Name,
				},
			},
			ChannelPrices: map[string]transfer.ChannelPrice{
				"default": {
					Price: high.Price,
				},
			},
		},
	}
}

func getProductVariantEntireFromHighliteMock(high highlite.Product) *transfer.VariantEntire {
	variant := getProductVariantFromHighliteMock(high).VariantEntire
	variant.ChannelPrices["default"] = transfer.ChannelPrice{
		Price: 12345,
	}

	return &variant
}
