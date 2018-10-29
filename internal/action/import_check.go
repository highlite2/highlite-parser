package action

import (
	"context"
	"os"
	"time"

	"highlite2-import/internal"
	"highlite2-import/internal/csv"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/log"
	"highlite2-import/internal/sylius/transfer"
)

// ImportCheck checks translations file
func ImportCheck(ctx context.Context, config internal.Config, logger log.ILogger) {
	ctx, cancel := context.WithTimeout(ctx, config.ImportTimeout)
	defer cancel()
	defer timeTrack(logger, time.Now(), "Check")

	file, err := os.Open(config.TranslationsFilePath)
	if err != nil {
		logger.Errorf("cant open translations file: %s", err)
		return
	}
	defer file.Close()

	dictionary, err := getHighliteTranslationsDictionary(config, file)
	if err != nil {
		logger.Errorf("cant fill dictionary: %s", err)
		return
	}

	for locale, translations := range dictionary.GetMap() {
		logger.Infof("there are %d translations in %s dictionary", len(translations), locale)
	}

	csvParser, err := getHighliteProductUpdatesCSVParser(ctx, config, logger)
	if err != nil {
		logger.Errorf("cant get highlite items reader: %s", err)
		return
	}

	csvMapper := csv.NewTitleMap(csvParser.GetNext())
	if err := csvMapper.CheckTitles(highlite.GetRequiredCSVTitles()); err != nil {
		logger.Errorf(err.Error())
		return
	}

	var productCounter, translationsMissing, everythingMissing int

	for csvParser.Next() {
		productCounter++
		product := highlite.GetProductFromCSVImport(csvMapper, csvParser.Values())
		if err := csvMapper.CheckValues(csvParser.Values()); err != nil {
			logger.Errorf("check values on line %d: %s", csvParser.CurrentRowIndex(), err)
			logger.Warn(product.String())
		}

		if translation, ok := dictionary.Get(transfer.LocaleRu, product.No); !ok || translation.Empty() {
			translationsMissing++
			logger.Warnf("missing translations for %s product", product.No)
			desc := product.GetDescription()
			short := product.GetShortDescription()
			if desc == "" && short == "" {
				everythingMissing++
				logger.Warnf("missing english version for http://highlite-spb.ru/products/%s", product.URL)
			}
		}
	}

	if csvParser.Err() != nil {
		logger.Errorf("error processing csv with product updates: %s", csvParser.Err())
	}

	logger.Infof("processed %d products", productCounter)
	if translationsMissing > 0 {
		logger.Warnf("missing translations for %d products", translationsMissing)
	}
	if everythingMissing > 0 {
		logger.Warnf("missing english version for %d products", everythingMissing)
	}

}
