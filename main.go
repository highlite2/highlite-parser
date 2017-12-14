package main

import (
	"context"
	"fmt"
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

	writer := internal.NewWriter(client, memo)

	file, err := os.Open("./_tmp/products_v1_0.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	parser := csv.NewReader(transform.NewReader(file, charmap.Windows1257.NewDecoder()))
	parser.ReadTitles()
	mapper := csv.NewTitleMap(parser.Titles())

	for {
		if !parser.Next() {
			break
		}

		product := highlite.GetProductFromCSVImport(mapper, parser.Values())
		fmt.Printf("%#v \n\n", product)

		writer.WriteProduct(ctx, product)

		break
	}

	if parser.Err() != nil {
		fmt.Println(parser.Err())
	}
}
