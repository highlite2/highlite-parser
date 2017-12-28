package main

import (
	"context"
	"time"

	"highlite-parser/internal"
	"highlite-parser/internal/cache"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/imprt"
	"highlite-parser/internal/log"
	"highlite-parser/internal/sylius"
)

func main() {
	config := internal.GetConfigFromFile("config/config.toml")

	logger := log.GetDefaultLog()

	highClient := highlite.NewClient(
		logger,
		config.Highlite.Login,
		config.Highlite.Password,
		config.Highlite.LoginEndpoint,
		config.Highlite.ItemsEndpoint,
	)

	client := sylius.NewClient(logger, config.Sylius.APIEndpoint, sylius.Auth{
		ClientID:     config.Sylius.ClientID,
		ClientSecret: config.Sylius.ClientSecret,
		Username:     config.Sylius.Username,
		Password:     config.Sylius.Password,
	})

	memo := cache.NewMemo()
	productImport := imprt.NewProductImport(client, memo, logger)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	processor := imprt.NewProcessor(logger, productImport, highClient)
	processor.Update(ctx)
}
