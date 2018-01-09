package highlite

import (
	"io"
	"fmt"
	"os"
	"highlite-parser/internal/csv"
)

const (
	// LangRU ..
	LangRU = "RU"

	translationOrderCode   = "order_code"
	translationCatTextMain = "cat_text_main"
	translationCatTextSubH = "cat_text_sub_h"
	translationUSP         = "usp"
	translationTechSpec    = "tech_spec"
)

var translationTitles map[string]string = map[string]string{
	LangRU: {
		translationOrderCode:   "Ordercode",
		translationCatTextMain: "cattext_main_rus",
		translationCatTextSubH: "cattext_subh_rus",
		translationUSP:         "USP_rus",
		translationTechSpec:    "techspec_rus",
	},
}

// Translation ... TODO
type Translation struct {
	lang   string
	reader io.Reader
	items  map[string]TranslationItem
}

// Get ... TODO
func (t *Translation) Get(id string) (TranslationItem, bool) {
	translation, ok := t.items[id]

	return translation, ok
}

// TranslationItem ... TODO
type TranslationItem struct {
	CatTextMain string
	CatTextSubH string
	USP         string
	TechSpec    string
}

// ProductDescription ... TODO
func (t *TranslationItem) ProductDescription() string {
	description := ""
	description += replaceHTMLEntities(t.USP)
	description += "\n"
	description += replaceHTMLEntities(t.CatTextMain)
	description += "\n\n"
	description += replaceHTMLEntities(t.TechSpec)

	return ""
}

func GetTranslationFromCSVFile(lang string, filePath string) (*Translation, error) {
	file, err := os.Open(filePath);
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvParser := csv.NewReader(file)
	csvParser.ReadTitles()
	csvMapper := csv.NewTitleMap(csvParser.Titles())

	translation := &Translation{
		lang: lang,
	}

	for ; csvParser.Next() ; {
		values := csvParser.Values()
		item := TranslationItem{
			CatTextMain: csvMapper.GetString(titleProductNo, values),
		}


	}

	return translation, nil
}