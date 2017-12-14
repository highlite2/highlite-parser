package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"highlite-parser/internal/csv"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func testGetTaxon(cl sylius.IClient, l log.ILogger) {
	ctx := context.Background()
	taxon, err := cl.GetTaxon(ctx, "category")
	if err != nil {
		l.Error(err.Error())
	} else {
		j, err := json.Marshal(taxon)
		if err != nil {
			l.Error(err.Error())
		} else {
			l.Info(string(j))
		}
	}
}

func main() {
	logger := log.GetDefaultLog()

	sylClient := sylius.NewClient(logger, "http://localhost:1221/app_dev.php/api", sylius.Auth{
		ClientID:     "demo_client",
		ClientSecret: "secret_demo_client",
		Username:     "api@example.com",
		Password:     "sylius-api",
	})

	testGetTaxon(sylClient, logger)

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

		p := highlite.GetProductFromCSVImport(mapper, parser.Values())
		fmt.Printf("%s \n", p.Category3.GetURL())
	}

	if parser.Err() != nil {
		fmt.Println(parser.Err())
	}
}
