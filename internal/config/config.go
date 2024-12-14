package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey    string
	SecretKey string
	BaseURL   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	env := os.Getenv("BINANCE_ENV")

	if env == "LIVE" {
		return &Config{
			APIKey:    os.Getenv("LIVE_BINANCE_API_KEY"),
			SecretKey: os.Getenv("LIVE_BINANCE_SECRET_KEY"),
			BaseURL:   os.Getenv("LIVE_BINANCE_BASE_URL"),
		}, nil
	}

	return &Config{
		APIKey:    os.Getenv("TEST_BINANCE_API_KEY"),
		SecretKey: os.Getenv("TEST_BINANCE_SECRET_KEY"),
		BaseURL:   os.Getenv("TEST_BINANCE_BASE_URL"),
	}, nil
}
