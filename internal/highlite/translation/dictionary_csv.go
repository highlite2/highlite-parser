package translation

import (
	"fmt"
	"os"

	"highlite-parser/internal/csv"
	"highlite-parser/internal/highlite"
)

const (
		translationTitleOrderCode =  "Ordercode"
		translationTitleCatTextMain = "cattext_main_rus"
		translationTitleCatTextSubH =  "cattext_subh_rus"
		translationTitleUSP =         "USP_rus"
		translationTitleTechSpec =    "techspec_rus"
)

// Translation ... TODO
type ItemCSV struct {
	CatTextMain string
	CatTextSubH string
	USP         string
	TechSpec    string
}

// GetDescription return product description from several fields
func (t *ItemCSV) GetDescription() string {
	description := ""
	description += highlite.ReplaceHTMLEntities(t.USP)
	description += "\n"
	description += highlite.ReplaceHTMLEntities(t.CatTextMain)
	description += "\n\n"
	description += highlite.ReplaceHTMLEntities(t.TechSpec)

	return description
}



// GetTranslationFromCSVFile ... TODO
func FillMemoryDictionaryFromRUCSV(dic *MemoryDictionary, langCode string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	csvParser := csv.NewReader(file)
	csvParser.SetSeparator(",")
	csvParser.ReadTitles()
	csvMapper := csv.NewTitleMap(csvParser.Titles())

	if _, ok := dic.languages[langCode]; !ok {
		dic.languages[langCode] = make(map[string]IItem)
	}

	for csvParser.Next() {
		values := csvParser.Values()

		productNo := csvMapper.GetString(translationTitleOrderCode, values)
		if productNo == "" {
			return fmt.Errorf("one of product codes is empty")
		}

		dic.languages[langCode][productNo] = &ItemCSV{
			CatTextMain: csvMapper.GetString(translationTitleCatTextMain, values),
			CatTextSubH: csvMapper.GetString(translationTitleCatTextSubH, values),
			USP:         csvMapper.GetString(translationTitleUSP, values),
			TechSpec:    csvMapper.GetString(translationTitleTechSpec, values),
		}
	}

	return nil
}
