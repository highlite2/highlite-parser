package action

import (
	"context"
	"fmt"
	"os"

	"highlite2-import/internal"
	"highlite2-import/internal/csv"
	"highlite2-import/internal/highlite/translation"
	"highlite2-import/internal/log"
	"highlite2-import/internal/sylius/transfer"
)

func TranslationCheck(ctx context.Context, config internal.Config, logger log.ILogger) {
	dictionary := translation.NewMemoryDictionary()

	file, err := os.Open(config.TranslationsFilePath)
	if err != nil {
		fmt.Printf("Can't open translations file: %s", err)
	}
	defer file.Close()

	csvParser := csv.NewReader(file)
	csvParser.QuotedQuotes = true
	csvParser.Separator = config.TranslationsFileSeparator

	if err := translation.FillMemoryDictionaryFromCSV(csvParser, dictionary, transfer.LocaleRu,
		translation.GetRussianTranslationsCSVTitles()); err != nil {
		fmt.Printf("Can't fill dictionary: %s", err.Error())
	} else {
		for locale, translations := range dictionary.GetMap() {
			fmt.Printf("There are %d translations in %s dictionary\n", len(translations), locale)
		}
	}
}
