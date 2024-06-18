package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/logger"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"60s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"60s"`
}

type Config struct {
	LogLevel    logger.LogLevel  `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat   logger.LogFormat `env:"LOG_FORMAT" envDefault:"text"`
	ApiKey      string           `env:"PORKBUN_API_KEY"`
	ApiSecret   string           `env:"PORKBUN_API_SECRET"`
	Health      ServerConfig     `envPrefix:"HEALTH"`
	HealthPort  int              `env:"HEALTH_PORT" envDefault:"8080"`
	Webhook     ServerConfig     `envPrefix:"WEBHOOK"`
	WebhookPort int              `env:"WEBHOOK_PORT" envDefault:"8888"`
}

func GetConfig() (*Config, error) {
	godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
