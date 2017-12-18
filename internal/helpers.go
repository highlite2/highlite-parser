package internal

import (
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/sylius/transfer"
)

// CreateNewTaxonFromHighliteCategory converts highlite category to sylius taxon struct
func CreateNewTaxonFromHighliteCategory(cat *highlite.Category) transfer.TaxonNew {
	taxon := transfer.TaxonNew{
		Code: cat.GetCode(),
		Translations: map[string]transfer.Translation{ // TODO take info from available locales from config
			transfer.LocaleEn: {
				Name: cat.Name,
				Slug: cat.GetURL(),
			},
			transfer.LocaleRu: {
				Name: cat.Name,
				Slug: cat.GetURL(),
			},
		},
	}

	if cat.Parent != nil {
		taxon.Parent = cat.Parent.GetCode()
	}

	return taxon
}

// CreateNewProductFromHighliteProduct converts highlite product to sylius product struct
func CreateNewProductFromHighliteProduct(p highlite.Product) transfer.ProductNew {
	product := transfer.ProductNew{
		Code: p.Code,
		Translations: map[string]transfer.Translation{ // TODO take info from available locales from config
			transfer.LocaleEn: {
				Name: p.Name,
				Slug: p.URL,
			},
			transfer.LocaleRu: {
				Name: p.Name,
				Slug: p.URL,
			},
		},
		MainTaxon:     p.Category3.GetCode(),
		ProductTaxons: p.Category3.GetCode(),
		Channels: []string{
			"US_WEB", // TODO take it from config
		},
	}

	return product
}
