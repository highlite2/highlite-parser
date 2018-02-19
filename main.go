package main

import (
	"context"
	"fmt"
	"flag"

	"highlite-parser/internal"
	"highlite-parser/internal/action"
	"highlite-parser/internal/log"
)

func main() {
	ctx := context.Background()
	config := internal.GetConfigFromFile("config/config.toml")
	logger := log.GetDefaultLog(config.LogLevel)

	act := flag.String("action", "", "Command")
	flag.Parse()

	switch *act {
	case "import":
		act := &action.HighliteImport{}
		act.Do(ctx, config, logger)
	case "tr":
		act := &action.CategoryTranslationTemplate{}
		err := act.Do(ctx, config, logger)
		if err != nil {
			logger.Error(err.Error())
		}
	default:
		fmt.Println("Please, specify a valid command.")
	}
}
