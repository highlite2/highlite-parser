package main

import (
	"context"
	"io"
	"os"
	"time"

	"highlite-parser/internal"
	"highlite-parser/internal/cache"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/highlite/translation"
	"highlite-parser/internal/imprt"
	"highlite-parser/internal/log"
	"highlite-parser/internal/queue"
	"highlite-parser/internal/sylius"
	"highlite-parser/internal/sylius/transfer"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	config := internal.GetConfigFromFile("config/config.toml")

	logger := log.GetDefaultLog(config.LogLevel)

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
	productImport := imprt.NewProductImport(syliusClient, memo, logger, dictionary)

	var itemsReader io.Reader

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

	jobPool := queue.NewPool(13)

	if itemsReader != nil {
		processor := imprt.NewProcessor(logger, jobPool, productImport, itemsReader)
		processor.Update(ctx)
	} else {
		logger.Error("Items reader is empty")
	}

	<-jobPool.Stop()
}
