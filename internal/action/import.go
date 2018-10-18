package action

import (
	"context"
	"io"
	"os"
	"time"

	"highlite2-import/internal"
	"highlite2-import/internal/cache"
	"highlite2-import/internal/csv"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/highlite/image"
	"highlite2-import/internal/highlite/translation"
	"highlite2-import/internal/imprt"
	"highlite2-import/internal/log"
	"highlite2-import/internal/queue"
	"highlite2-import/internal/sylius"
	"highlite2-import/internal/sylius/transfer"
)

// Import uploads products
func Import(ctx context.Context, config internal.Config, logger log.ILogger) {
	ctx, cancel := context.WithTimeout(ctx, config.ImportTimeout)
	defer cancel()
	defer timeTrack(logger, time.Now(), "Import")

	translationsFile, err := os.Open(config.TranslationsFilePath)
	if err != nil {
		logger.Errorf("cant open translations translationsFile: %s", err)
		return
	}
	defer translationsFile.Close()

	dictionary, err := getHighliteTranslationsDictionary(config, translationsFile)
	if err != nil {
		logger.Errorf("cant fill dictionary: %s", err)
		return
	}

	syliusClient := getSyliusClient(config, logger)
	memo := cache.NewMemo()
	categoryImport := imprt.NewCategoryImport(syliusClient, memo, logger)
	productImport := imprt.NewProductImport(syliusClient, categoryImport, logger, dictionary, image.HTTPProvider{})

	csvParser, err := getHighliteProductUpdatesCSVParser(ctx, config, logger)
	if err != nil {
		logger.Errorf("can't get highlite items reader: %s", err)
		return
	}

	jobPool := queue.NewPool(10)
	processor := imprt.NewProcessor(logger, jobPool, productImport, csvParser)
	processor.Update(ctx)

	<-jobPool.Stop()
}

func getSyliusClient(config internal.Config, logger log.ILogger) *sylius.Client {
	return sylius.NewClient(logger, config.Sylius.APIEndpoint, sylius.Auth{
		ClientID:     config.Sylius.ClientID,
		ClientSecret: config.Sylius.ClientSecret,
		Username:     config.Sylius.Username,
		Password:     config.Sylius.Password,
	})
}

func getHighliteTranslationsDictionary(config internal.Config, items io.Reader) (*translation.MemoryDictionary, error) {
	csvParser := csv.NewReader(items)
	csvParser.QuotedQuotes = true
	csvParser.Separator = config.TranslationsFileSeparator

	dictionary := translation.NewMemoryDictionary()
	if err := translation.FillMemoryDictionaryFromCSV(
		csvParser,
		dictionary,
		transfer.LocaleRu,
		translation.GetRussianTranslationsCSVTitles(),
	); err != nil {
		return nil, err
	}

	return dictionary, nil
}

func getHighliteProductUpdatesCSVParser(ctx context.Context, config internal.Config, logger log.ILogger) (*csv.Reader, error) {
	highClient := highlite.NewClient(
		logger,
		config.Highlite.Login,
		config.Highlite.Password,
		config.Highlite.LoginEndpoint,
		config.Highlite.ItemsEndpoint,
	)

	reader, err := highClient.GetItemsReader(ctx)
	if err != nil {
		return nil, err
	}

	csvParser := csv.NewReader(reader)
	csvParser.Separator = ';'
	csvParser.FieldsFixed = false
	csvParser.OneRowRecord = true

	return csvParser, nil
}

func timeTrack(logger log.ILogger, start time.Time, name string) {
	elapsed := time.Since(start)
	logger.Infof("[%s] took %s", name, elapsed)
}
