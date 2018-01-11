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

// GetTranslationFromCSVFile ... TODO
func GetTranslationFromCSVFile(lang string, filePath string) (*DictionaryMap, error) {
	if err := checkTranslationTitlesForLanguage(lang); err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvParser := csv.NewReader(file)
	csvParser.SetSeparator(",")
	csvParser.ReadTitles()
	csvMapper := csv.NewTitleMap(csvParser.Titles())

	translation := &DictionaryMap{
		translations: map[string]map[string]Translation{
			lang: {},
		},
	}

	titles := csvTranslationsTitles[lang]

	for csvParser.Next() {
		values := csvParser.Values()

		productNo := csvMapper.GetString(titles[translationTitleOrderCode], values)
		if productNo == "" {
			return nil, fmt.Errorf("one of product codes is empty")
		}

		translation.translations[lang][productNo] = Translation{
			CatTextMain: csvMapper.GetString(titles[translationTitleCatTextMain], values),
			CatTextSubH: csvMapper.GetString(titles[translationTitleCatTextSubH], values),
			USP:         csvMapper.GetString(titles[translationTitleUSP], values),
			TechSpec:    csvMapper.GetString(titles[translationTitleTechSpec], values),
		}
	}

	return translation, nil
}
