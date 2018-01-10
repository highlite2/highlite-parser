package highlite

import (
	"fmt"
	"os"

	"highlite-parser/internal/csv"
)

const (
	// LangRU ...
	LangRU = "RU"

	translationTitleOrderCode = iota
	translationTitleCatTextMain
	translationTitleCatTextSubH
	translationTitleUSP
	translationTitleTechSpec
)

var csvTranslationsTitles = map[string]map[int]string{
	LangRU: {
		translationTitleOrderCode:   "Ordercode",
		translationTitleCatTextMain: "cattext_main_rus",
		translationTitleCatTextSubH: "cattext_subh_rus",
		translationTitleUSP:         "USP_rus",
		translationTitleTechSpec:    "techspec_rus",
	},
}

func checkTranslationTitlesForLanguage(lang string) error {
	if _, ok := csvTranslationsTitles[lang]; !ok {
		return fmt.Errorf("undefined lang")
	}

	checkList := []int{
		translationTitleOrderCode,
		translationTitleCatTextMain,
		translationTitleCatTextSubH,
		translationTitleUSP,
		translationTitleTechSpec,
	}

	for _, key := range checkList {
		if _, ok := csvTranslationsTitles[lang][key]; !ok {
			return fmt.Errorf("key [%s] is undefined for %s lang", key, lang)
		}
	}

	return nil
}

// Translation ... TODO
type Translation struct {
	lang  string
	items map[string]TranslationItem
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

// GetTranslationFromCSVFile ... TODO
func GetTranslationFromCSVFile(lang string, filePath string) (*Translation, error) {
	if err := checkTranslationTitlesForLanguage(lang); err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvParser := csv.NewReader(file)
	csvParser.ReadTitles()
	csvMapper := csv.NewTitleMap(csvParser.Titles())

	translation := &Translation{
		lang:  lang,
		items: make(map[string]TranslationItem),
	}

	titles := csvTranslationsTitles[lang]

	for csvParser.Next() {
		values := csvParser.Values()

		productNo := csvMapper.GetString(titles[translationTitleOrderCode], values)
		if productNo == "" {
			return nil, fmt.Errorf("one of product codes is empty")
		}

		translation.items[productNo] = TranslationItem{
			CatTextMain: csvMapper.GetString(titles[translationTitleCatTextMain], values),
			CatTextSubH: csvMapper.GetString(titles[translationTitleCatTextSubH], values),
			USP:         csvMapper.GetString(titles[translationTitleUSP], values),
			TechSpec:    csvMapper.GetString(titles[translationTitleTechSpec], values),
		}
	}

	return translation, nil
}
