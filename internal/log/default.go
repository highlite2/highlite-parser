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
	case "warn":
		apexLog.SetLevel(apexLog.WarnLevel)
	case "error":
		apexLog.SetLevel(apexLog.ErrorLevel)
	case "fatal":
		apexLog.SetLevel(apexLog.FatalLevel)
	default:
		apexLog.SetLevel(apexLog.FatalLevel)
		apexLog.Log.Warnf("Unknown log level type: %s, setting level to fatal level", level)
	}

	return apexLog.Log
}
