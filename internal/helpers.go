package internal

import (
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/sylius/transfer"
)

// CreateNewTaxonFromHighliteCategory converts highlite category to taxon struct
func CreateNewTaxonFromHighliteCategory(cat *highlite.Category) transfer.TaxonNew {
	taxon := transfer.TaxonNew{
		Code: cat.GetCode(),
		Translations: map[string]transfer.Translation{
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
