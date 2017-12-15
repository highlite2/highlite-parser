package main

import (
	"context"
	"os"

	"highlite-parser/internal"
	"highlite-parser/internal/cache"
	"highlite-parser/internal/csv"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {
	ctx := context.Background()

	logger := log.GetDefaultLog()

	client := sylius.NewClient(logger, "http://localhost:1221/app_dev.php/api", sylius.Auth{
		ClientID:     "demo_client",
		ClientSecret: "secret_demo_client",
		Username:     "api@example.com",
		Password:     "sylius-api",
	})

	memo := cache.NewMemo()

	writer := internal.NewWriter(client, memo, logger)

	file, err := os.Open("./_tmp/products_v1_0.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	parser := csv.NewReader(transform.NewReader(file, charmap.Windows1257.NewDecoder()))
	parser.ReadTitles()
	mapper := csv.NewTitleMap(parser.Titles())

	logger.Info("Start csv file processing")

	for {
		if !parser.Next() {
			break
		}

		pr := highlite.GetProductFromCSVImport(mapper, parser.Values())
		logger.Debugf("Processing pr: %s", pr.Category3.GetURL())

		if err := writer.WriteProduct(ctx, pr); err != nil {
			logger.Errorf("Product processing error: %s", err.Error())
		}
	}

	if parser.Err() != nil {
		logger.Errorf("Csv processing error: %s", parser.Err().Error())
	}

	logger.Info("Stop csv file processing")
}
