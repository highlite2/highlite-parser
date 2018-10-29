package translation

import (
	"fmt"
	"strings"

	"highlite2-import/internal/csv"
)

const (
	titleOrderCode = iota
	titleCatTextMain
	titleCatTextSubH
	titleUSP
	titleTechSpec
)

// GetRussianTranslationsCSVTitles returns titles for Russian translations csv file.
func GetRussianTranslationsCSVTitles() map[int]string {
	return map[int]string{
		titleOrderCode:   "Ordercode",
		titleCatTextMain: "cattext_main_rus",
		titleCatTextSubH: "cattext_subh_rus",
		titleUSP:         "USP_rus",
		titleTechSpec:    "techspec_rus",
	}
}

// ProductCSV is an implementation of IProduct.
type ProductCSV struct {
	MainText   string
	SubHeading string
	USP        string
	TechSpec   string
}

// GetDescription returns a product description.
func (t *ProductCSV) GetDescription() string {
	description := ""
	if t.USP != "" {
		description += t.USP
		description += "\n"
	}

	description += t.MainText
	description += "\n\n"
	description += t.TechSpec

	return strings.Trim(description, "\n")
}

// GetShortDescription returns short product description.
func (t *ProductCSV) GetShortDescription() string {
	return strings.Trim(t.SubHeading, "\n")
}

// Empty returns true if all translations are empty strings
func (t *ProductCSV) Empty() bool {
	return t.GetDescription() == "" && t.GetShortDescription() == ""
}

// FillMemoryDictionaryFromCSV fills an exact MemoryDictionary with translations from a csv file.
func FillMemoryDictionaryFromCSV(csvParser *csv.Reader, dic *MemoryDictionary, lang string, titles map[int]string) error {
	csvMapper := csv.NewTitleMap(csvParser.GetNext())

	if _, ok := dic.languages[lang]; !ok {
		dic.languages[lang] = make(map[string]IProduct)
	}

	for csvParser.Next() {
		values := csvParser.Values()

		productNo := csvMapper.GetString(titles[titleOrderCode], values)
		if productNo == "" {
			return fmt.Errorf("one of product codes is empty")
		}

		dic.languages[lang][productNo] = &ProductCSV{
			MainText:   csvMapper.GetString(titles[titleCatTextMain], values),
			SubHeading: csvMapper.GetString(titles[titleCatTextSubH], values),
			USP:        csvMapper.GetString(titles[titleUSP], values),
			TechSpec:   csvMapper.GetString(titles[titleTechSpec], values),
		}
	}

	return csvParser.Err()
}
