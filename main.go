package main

import (
	"context"
	"encoding/json"

	apexLog "github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"os"

	"fmt"

	"highlite-parser/internal"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/sylius"
)

func getLog() internal.ILogger {
	apexLog.SetHandler(cli.Default)
	apexLog.SetLevel(apexLog.DebugLevel)
	return apexLog.Log
}

func testGetTaxon(cl sylius.IClient, l internal.ILogger) {
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
	log := getLog()

	syliusClient := sylius.NewClient(log, "http://localhost:1221/app_dev.php/api", sylius.Auth{
		ClientID:     "demo_client",
		ClientSecret: "secret_demo_client",
		Username:     "api@example.com",
		Password:     "sylius-api",
	})

	testGetTaxon(syliusClient, log)

	fmt.Println()
	fmt.Println()

	file, err := os.Open("./_tmp/products_v1_0.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	parser := highlite.NewCSVReader(highlite.GetWindows1257Decoder(file), log)
	parser.ReadTitles()
	for {
		if !parser.Next() {
			break
		}
		m := parser.TitledValues()
		p := highlite.GetProductFromCSVImport(m)
		fmt.Printf("%s \n", p.Category3.GetSlug())
	}
}
