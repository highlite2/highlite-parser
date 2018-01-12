package translation

import (
	"fmt"
	"os"

	"highlite-parser/internal/csv"
	"highlite-parser/internal/highlite"
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
	description += highlite.ReplaceHTMLEntities(t.USP)
	description += "\n"
	description += highlite.ReplaceHTMLEntities(t.MainText)
	description += "\n\n"
	description += highlite.ReplaceHTMLEntities(t.TechSpec)

	return description
}

// GetShortDescription returns short product description.
func (t *ProductCSV) GetShortDescription() string {
	return highlite.ReplaceHTMLEntities(t.SubHeading)
}

// FillMemoryDictionaryFromCSV fills an exact MemoryDictionary with translations from a csv file.
func FillMemoryDictionaryFromCSV(dic *MemoryDictionary, lang string, path string, titles map[int]string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	csvParser := csv.NewReader(file)
	csvParser.SetSeparator(",")
	csvParser.ReadTitles()
	csvMapper := csv.NewTitleMap(csvParser.Titles())

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

	return nil
}
