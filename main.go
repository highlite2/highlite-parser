package main

import (
	"time"

	apexLog "github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"context"
	"encoding/json"

	"highlite-parser/internal"
	"highlite-parser/internal/client/sylius"
)

func getLog() internal.ILogger {
	apexLog.SetHandler(cli.Default)
	apexLog.SetLevel(apexLog.DebugLevel)
	return apexLog.Log
}

func main() {
	l := getLog()
	cl := sylius.NewClient(l, "http://localhost:1221/app_dev.php/api", sylius.Auth{
		ClientID:     "demo_client",
		ClientSecret: "secret_demo_client",
		Username:     "api@example.com",
		Password:     "sylius-api",
	})

	ctx := context.Background()
	taxon, err := cl.GetTaxon(ctx, "t_shirts")
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

	time.Sleep(time.Second * 5)
}
