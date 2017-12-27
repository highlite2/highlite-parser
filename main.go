package main

import (
	"context"
	"os"
	"time"

	"highlite-parser/internal"
	"highlite-parser/internal/cache"
	"highlite-parser/internal/csv"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"

	"github.com/go-resty/resty"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := internal.GetConfigFromFile("config/config.toml")

	logger := log.GetDefaultLog()

	resty.DefaultClient.Debug = false
	client := sylius.NewClient(logger, cfg.Sylius.APIEndpoint, sylius.Auth{
		ClientID:     cfg.Sylius.ClientID,
		ClientSecret: cfg.Sylius.ClientSecret,
		Username:     cfg.Sylius.Username,
		Password:     cfg.Sylius.Password,
	})

	memo := cache.NewMemo()

	productImport := internal.NewProductImport(client, memo, logger)

	file, err := os.Open("./_tmp/products_v1_0.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	parser := csv.NewReader(transform.NewReader(file, charmap.Windows1257.NewDecoder()))
	parser.ReadTitles()
	mapper := csv.NewTitleMap(parser.Titles())

	logger.Info("Start csv file processing")

	i := 0
	for run := true; run && i < 1; i++ {
		select {

		case <-ctx.Done():
			logger.Info("Context timeout")
			run = false

		default:
			if !parser.Next() {
				run = false
				break
			}

			pr := highlite.GetProductFromCSVImport(mapper, parser.Values())
			logger.Debugf("Processing product: category %s", pr.Category3.GetURL())

			if err := productImport.Import(ctx, pr); err != nil {
				logger.Errorf("Product processing error: %s", err.Error())
			}

		}
	}

	if parser.Err() != nil {
		logger.Errorf("Csv processing error: %s", parser.Err().Error())
	}

	logger.Info("Stop csv file processing")
}
