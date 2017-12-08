package main

import (
	"highlite-parser/client/sylius"
	"highlite-parser/log"
	"time"

	apexLog "github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func getLog() log.Logger {
	apexLog.SetHandler(cli.Default)
	apexLog.SetLevel(apexLog.DebugLevel)
	return apexLog.Log
}

func main() {
	logger := getLog()
	_ = sylius.NewClient(logger, "http://localhost:1221/app_dev.php/api", sylius.Auth{
		ClientID:     "3u721kcbho4kcosgws08s84gw48wc0g40ggc088s8ccs8s40w0",
		ClientSecret: "tplxj5h4e800gc8480ckss0okc8kwccck4ks4o40ckoc0c4w",
		Username:     "test@test.com",
		Password:     "123123",
	})

	time.Sleep(time.Second * 30)
}

/*
curl http://localhost:1221/app_dev.php/api/oauth/v2/token \
    -d "client_id"=3u721kcbho4kcosgws08s84gw48wc0g40ggc088s8ccs8s40w0 \
    -d "client_secret"=tplxj5h4e800gc8480ckss0okc8kwccck4ks4o40ckoc0c4w \
    -d "grant_type"=password \
    -d "username"=test@test.com \
    -d "password"=123123


curl http://localhost:1221/app_dev.php/api/v1/taxons/toys \
    -H "Authorization: Bearer OTkzNjE4M2I2YWYxNWM3MDA4MTdmNmUyYjIwZTcyN2Y3ZjNhNjRlMjc2ZWI3OTA0MjE5NmI2ZjBiMzc5MzMyZQ" \
    -H "Accept: application/json"

curl http://localhost:1221/app_dev.php/api/v1/users/ \
    -H "Authorization: Bearer OTkzNjE4M2I2YWYxNWM3MDA4MTdmNmUyYjIwZTcyN2Y3ZjNhNjRlMjc2ZWI3OTA0MjE5NmI2ZjBiMzc5MzMyZQ"

curl http://localhost:1221/app_dev.php/api/v1/users/3 \
    -H "Authorization: Bearer M2EyM2JjZjlkYzdjZTc0NGNiNGYxNDQ2Y2NiMmMyMjk0MTg1MjQzMWExMjk3NjhmOTE5YmI5ZGY1NmVhMmE5NQ"
*/
