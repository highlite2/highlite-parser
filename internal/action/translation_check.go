package action

import (
	"context"
	"os"

	"highlite2-import/internal"
	"highlite2-import/internal/log"
)

// TranslationCheck checks translations file
func TranslationCheck(ctx context.Context, config internal.Config, logger log.ILogger) {
	file, err := os.Open(config.TranslationsFilePath)
	if err != nil {
		logger.Errorf("cant open translations file: %s", err)
		return
	}
	defer file.Close()

	dictionary, err := getHighliteTranslationsDictionary(config, file)
	if err != nil {
		logger.Errorf("cant fill dictionary: %s", err)
	} else {
		for locale, translations := range dictionary.GetMap() {
			logger.Infof("there are %d translations in %s dictionary", len(translations), locale)
		}
	}
}
