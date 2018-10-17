package action

import (
	"highlite2-import/internal"
	"highlite2-import/internal/log"
	"context"
	"highlite2-import/internal/highlite/translation"
	"highlite2-import/internal/sylius/transfer"
	"fmt"
)

func TranslationCheck(ctx context.Context, config internal.Config, logger log.ILogger) {
	dictionary := translation.NewMemoryDictionary()
	if err := translation.FillMemoryDictionaryFromCSV(dictionary, transfer.LocaleRu,
		config.TranslationsFilePath, translation.GetRussianTranslationsCSVTitles()); err != nil {
		fmt.Printf("Can't fill dictionary: %s", err.Error())
	} else {
		for locale, translations := range dictionary.GetMap() {
			fmt.Printf("There are %d translations in %s dictionary\n", len(translations), locale)
		}
	}
}
