package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	ApiKey    string `env:"PORKBUN_API_KEY"`
	ApiSecret string `env:"PORKBUN_API_SECRET"`
}

func GetConfig() (*Config, error) {
	godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
