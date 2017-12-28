package log

import (
	apexLog "github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// GetDefaultLog returns log with debug level and cli output
func GetDefaultLog(level string) ILogger {
	apexLog.SetHandler(cli.Default)

	switch level {
	case "debug":
		apexLog.SetLevel(apexLog.DebugLevel)
	case "info":
		apexLog.SetLevel(apexLog.InfoLevel)
	default:
		apexLog.SetLevel(apexLog.FatalLevel)
	}

	return apexLog.Log
}
