package highlite

import (
	"strings"

	"highlite-parser/internal/csv"
)

const (
	// StatusCurrent ...
	StatusCurrent string = "CURRENT"
	// StatusDecline ...
	StatusDecline string = "DECLINE"
	// StatusSpecial ...
	StatusSpecial string = "SPECIAL"
	// StatusNew ...
	StatusNew string = "NEW"
	// StatusEOL ...
	StatusEOL string = "EOL"
	// StatusOnRequest ...
	StatusOnRequest string = "ONREQUEST"

	// InStockYes ...
	InStockYes string = "YES"
	// InStockNo ...
	InStockNo string = "NO"
	// InStockLow ...
	InStockLow string = "LOW"
)

const (
	titleProductNo         string = "Product No."
	titleProductName       string = "Product Name"
	titleCountryOfOrigin   string = "Country of Origin"
	titleWeight            string = "Weight (kg)"
	titleLength            string = "Length (m)"
	titleWidth             string = "Width (m)"
	titleHeight            string = "Height (m)"
	titleQuantityInCarton  string = "Qty. in Carton"
	titleEANNo             string = "EAN No."
	titleInternetAdvice    string = "Internet Advice Prijs"
	titleUnitPrice         string = "Unit Price"
	titleMinSalesQuantity  string = "Min. Sales Qty."
	titleCategory          string = "Category"
	titleSubcategory1      string = "Subcategory 1"
	titleSubcategory2      string = "Subcategory 2"
	titleTariffNo          string = "Tariff No."
	titleStatus            string = "Status"
	titleAccessory         string = "Accessory"
	titleSubstitute        string = "Substitute"
	titleCatalogPage       string = "Catalog Page"
	titleInStock           string = "In Stock"
	titleExpWeekOfArrival  string = "\"Exp. Week of Arrival\""
	titleWebshop           string = "Webshop"
	titleSubheadingEN      string = "Subheading EN"
	titleMainDescriptionEN string = "Main Description EN"
	titleSpecsEN           string = "Specs EN"
	titleImagesList        string = "Images List"
	titleManual            string = "Manual"
	titleYoutubeLink       string = "Youtube Link"
	titleBrand             string = "Brand"
	titleSubcategory3      string = "Subcategory 3"
)

// Product is a highlite product.
type Product struct {
	Identification

	No          string
	Name        string
	SubHeading  string
	Description string
	Specs       string
	Brand       string

	InStock string
	Status  string

	Country string
	Weight  float64
	Length  float64
	Width   float64
	Height  float64
	Price   float64

	Category1 *Category
	Category2 *Category
	Category3 *Category

	Images []string
}

// ProductDescription combines description and specs and removes html entities.
func (p Product) ProductDescription() string {
	description := ""
	description += ReplaceHTMLEntities(p.Description)
	description += "\n\n"
	description += ReplaceHTMLEntities(p.Specs)

	return description
}

// GetProductFromCSVImport creates product object from csv import data.
func GetProductFromCSVImport(mapper *csv.TitleMap, values []string) Product {
	cat1 := NewCategory(mapper.GetString(titleCategory, values), nil)
	cat2 := NewCategory(mapper.GetString(titleSubcategory1, values), cat1)
	cat3 := NewCategory(mapper.GetString(titleSubcategory2, values), cat2)

	product := Product{
		No:          mapper.GetString(titleProductNo, values),
		Name:        mapper.GetString(titleProductName, values),
		SubHeading:  mapper.GetString(titleSubheadingEN, values),
		Description: mapper.GetString(titleMainDescriptionEN, values),
		Specs:       mapper.GetString(titleSpecsEN, values),
		Brand:       mapper.GetString(titleBrand, values),

		InStock: mapper.GetString(titleInStock, values),
		Status:  mapper.GetString(titleStatus, values),

		Country: mapper.GetString(titleCountryOfOrigin, values),
		Weight:  mapper.GetFloat(titleWeight, values),
		Length:  mapper.GetFloat(titleLength, values),
		Width:   mapper.GetFloat(titleWidth, values),
		Height:  mapper.GetFloat(titleHeight, values),
		Price:   mapper.GetFloat(titleUnitPrice, values),

		Category1: cat1,
		Category2: cat2,
		Category3: cat3,
	}

	product.SetCodeAndURL(product.Name + " " + product.No)
	product.Code = "highlite-" + product.No

	imagesString := mapper.GetString(titleImagesList, values)
	product.Images = strings.Fields(strings.Replace(imagesString, "|", " ", -1))

	return product
}

// Removes specific for highlite html tags.
func ReplaceHTMLEntities(str string) string {
	return strings.Replace(str, "<br />", "\n", -1)
}
