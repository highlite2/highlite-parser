package main

import (
	"context"
	"flag"
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
	default:
		logger.Info("checking translations...")
		action.TranslationCheck(ctx, config, logger)

		logger.Info("checking product updates...")
		action.UpdatesCheck(ctx, config, logger)
	}
}
