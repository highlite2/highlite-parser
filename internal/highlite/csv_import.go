package highlite

const (
	titleProductNo         string = "Product No."
	titleProductName       string = "Product Name"
	titleCountryOfOrigin   string = "Country of Origin"
	titleWeight            string = "Weight (kg)"
	titleLength            string = "Length (m)"
	titleWidth             string = "Width (m)"
	titleHeight            string = "Height (m)"
	titleQuantity          string = "Qty. in Carton"
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
	titleBrand             string = "Brand"
	titleSubcategory3      string = "Subcategory 3"
)

// GetProductFromCSVImport creates product object from csv import data.
func GetProductFromCSVImport(input map[string]string) *Product {
	cat1 := NewCategory(input[titleCategory], nil)
	cat2 := NewCategory(input[titleSubcategory1], cat1)
	cat3 := NewCategory(input[titleSubcategory2], cat2)

	product := &Product{
		Category1: cat1,
		Category2: cat2,
		Category3: cat3,
	}

	return product
}
