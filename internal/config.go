package internal

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// GetConfigFromFile creates config instance and makes all initializations.
func GetConfigFromFile(path string) Config {
	var cfg = Config{}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		panic(fmt.Errorf("can't parse config: %s", err.Error()))
	}

	return cfg
}

// Config is an application global config.
type Config struct {
	LogLevel string

	Highlite struct {
		Login         string
		Password      string
		LoginEndpoint string
		ItemsEndpoint string
	}

	Sylius struct {
		ClientID     string
		ClientSecret string
		Username     string
		Password     string
		APIEndpoint  string
	}
}
