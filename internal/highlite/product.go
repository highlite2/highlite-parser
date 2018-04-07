package highlite

import (
	"strings"

	"highlite2-import/internal/csv"
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
	titleExpWeekOfArrival  string = "Exp. Week of Arrival"
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
	Price   float64

	CategoryRoot *Category
	Category1    *Category
	Category2    *Category
	Category3    *Category

	Images []string
}

// GetProductDescription combines description and specs and removes html entities.
func (p Product) GetProductDescription() string {
	description := replaceHTMLEntities(p.Description)
	description += "\n\n"
	description += replaceHTMLEntities(p.Specs)

	return strings.Trim(description, "\n ")
}

// GetShortProductDescription returns trimmed short product description
func (p Product) GetShortProductDescription() string {
	return strings.Trim(p.SubHeading, "\n ")
}

// GetProductName combines name from brand and name fields.
func (p Product) GetProductName() string {
	var name string
	if p.Brand == "" {
		name = p.Name
	} else {
		name = p.Brand + " " + p.Name
	}

	return strings.Trim(name, "\n ")
}

// GetBrandCode turns brand name into brand code.
func (p Product) GetBrandCode() string {
	brand := strings.ToLower(p.Brand)
	brand = categoryCodeRegExp.ReplaceAllString(brand, " ")
	fields := strings.Fields(brand)

	return strings.Join(fields, "-")
}

// GetProductFromCSVImport creates product object from csv import data.
func GetProductFromCSVImport(mapper *csv.TitleMap, values []string) Product {
	cat0 := NewCategory("Category", nil)
	cat0.Root = true
	cat1 := NewCategory(mapper.GetString(titleCategory, values), cat0)
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
		Price:   mapper.GetFloat(titleUnitPrice, values),

		CategoryRoot: cat0,
		Category1:    cat1,
		Category2:    cat2,
		Category3:    cat3,
	}

	product.SetCodeAndURL(product.Name + " " + product.No)
	product.Code = "highlite-" + product.No

	imagesString := mapper.GetString(titleImagesList, values)
	product.Images = parseImages(imagesString)

	return product
}

// Converts string containing product images to a slice of strings.
func parseImages(str string) []string {
	split := strings.Split(str, "|")
	images := make([]string, 0, len(split))
	for _, im := range split {
		im = strings.Trim(im, " ")
		if im != "" {
			images = append(images, im)
		}
	}

	return images
}

// Removes specific for highlite html tags.
func replaceHTMLEntities(str string) string {
	return strings.Replace(str, "<br />", "\n", -1)
}
