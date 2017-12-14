package log

import (
	apexLog "github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// GetDefaultLog returns log with debug level and cli output
func GetDefaultLog() ILogger {
	apexLog.SetHandler(cli.Default)
	apexLog.SetLevel(apexLog.DebugLevel)

	return apexLog.Log
}
