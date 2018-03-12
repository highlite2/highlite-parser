package action

import (
	"context"
	"io"
	"os"
	"time"

	"highlite2-import/internal"
	"highlite2-import/internal/cache"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/highlite/image"
	"highlite2-import/internal/highlite/translation"
	"highlite2-import/internal/imprt"
	"highlite2-import/internal/log"
	"highlite2-import/internal/queue"
	"highlite2-import/internal/sylius"
	"highlite2-import/internal/sylius/transfer"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type HighliteImport struct{}

func (h *HighliteImport) Do(ctx context.Context, config internal.Config, logger log.ILogger) {
	ctx, cancel := context.WithTimeout(ctx, time.Hour*3)
	defer cancel()

	defer timeTrack(logger, time.Now(), "Import")

	highClient := highlite.NewClient(
		logger,
		config.Highlite.Login,
		config.Highlite.Password,
		config.Highlite.LoginEndpoint,
		config.Highlite.ItemsEndpoint,
	)

	syliusClient := sylius.NewClient(logger, config.Sylius.APIEndpoint, sylius.Auth{
		ClientID:     config.Sylius.ClientID,
		ClientSecret: config.Sylius.ClientSecret,
		Username:     config.Sylius.Username,
		Password:     config.Sylius.Password,
	})

	dictionary := translation.NewMemoryDictionary()
	if err := translation.FillMemoryDictionaryFromCSV(dictionary, transfer.LocaleRu,
		config.TranslationsFilePath, translation.GetRussianTranslationsCSVTitles()); err != nil {
		logger.Errorf("Can't fill dictionary: %s", err.Error())

		return
	}

	memo := cache.NewMemo()
	categoryImport := imprt.NewCategoryImport(syliusClient, memo, logger)
	productImport := imprt.NewProductImport(syliusClient, categoryImport, logger, dictionary, image.HTTPProvider{})

	var itemsReader io.Reader
	// TODO refactor logic of reader creating
	if config.ItemsFilePath == "" {
		if reader, err := highClient.GetItemsReader(ctx); err != nil {
			logger.Errorf("Can't get highlite items reader: %s", err.Error())
		} else {
			itemsReader = reader
		}
	} else {
		if file, err := os.Open(config.ItemsFilePath); err != nil {
			logger.Errorf("Can't open file for reading items: %s", err.Error())
		} else {
			defer file.Close()
			itemsReader = transform.NewReader(file, charmap.Windows1257.NewDecoder())
		}
	}

	jobPool := queue.NewPool(10)

	if itemsReader != nil {
		processor := imprt.NewProcessor(logger, jobPool, productImport, itemsReader)
		processor.Update(ctx)
	} else {
		logger.Error("Items reader is empty")
	}

	<-jobPool.Stop()
}

// Time logging
func timeTrack(logger log.ILogger, start time.Time, name string) {
	elapsed := time.Since(start)
	logger.Infof("[%s] took %s", name, elapsed)
}
