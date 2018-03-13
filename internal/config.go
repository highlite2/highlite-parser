package internal

import "os"

// GetConfigFromEnv creates config instance and makes all initializations.
func GetConfigFromEnv() Config {
	var cfg = Config{}

	cfg.TranslationsFilePath = os.Getenv("TRANSLATIONS_FILE_PATH")

	cfg.LogLevel = os.Getenv("LOG_LEVEL")

	cfg.Highlite.Login = os.Getenv("HIGHLITE_LOGIN")
	cfg.Highlite.Password = os.Getenv("HIGHLITE_PASSWORD")
	cfg.Highlite.LoginEndpoint = os.Getenv("HIGHLITE_LOGIN_ENDPOINT")
	cfg.Highlite.ItemsEndpoint = os.Getenv("HIGHLITE_ITEMS_ENDPOINT")

	cfg.Sylius.ClientID = os.Getenv("SYLIUS_CLIENT_ID")
	cfg.Sylius.ClientSecret = os.Getenv("SYLIUS_CLIENT_SECRET")
	cfg.Sylius.Username = os.Getenv("SYLIUS_USERNAME")
	cfg.Sylius.Password = os.Getenv("SYLIUS_PASSWORD")
	cfg.Sylius.APIEndpoint = os.Getenv("SYLIUS_API_ENDPOINT")

	return cfg
}

// Config is an application global config.
type Config struct {
	TranslationsFilePath string

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
