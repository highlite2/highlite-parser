package main

import (
	"context"
	"flag"
	"fmt"

	"highlite2-import/internal"
	"highlite2-import/internal/action"
	"highlite2-import/internal/log"
)

func main() {
	ctx := context.Background()
	config := internal.GetConfigFromEnv()
	logger := log.GetDefaultLog(config.LogLevel)

	act := flag.String("action", "", "Command")
	flag.Parse()

	switch *act {
	case "import":
		action.Import(ctx, config, logger)
	case "translation-check":
		action.TranslationCheck(ctx, config, logger)
	default:
		fmt.Println("no cmd was specified")
	}
}
