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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	logger := log.GetDefaultLog()

	resty.DefaultClient.Debug = false
	client := sylius.NewClient(logger, "http://localhost:1221/app_dev.php/api", sylius.Auth{
		ClientID:     "demo_client",
		ClientSecret: "secret_demo_client",
		Username:     "api@example.com",
		Password:     "sylius-api",
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
