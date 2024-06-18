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

type PorkbunConfig struct {
	ApiKey    string `env:"API_KEY"`
	ApiSecret string `env:"API_SECRET"`
}

type ExternalDnsDomainConfig struct {
	DomainFilter         []string `env:"DOMAIN_FILTER" default:""`
	ExcludeDomains       []string `env:"EXCLUDE_DOMAIN_FILTER" default:""`
	RegexDomainFilter    string   `env:"REGEXP_DOMAIN_FILTER" default:""`
	RegexDomainExclusion string   `env:"REGEXP_DOMAIN_FILTER_EXCLUSION" default:""`
}

type Config struct {
	LogLevel     logger.LogLevel  `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat    logger.LogFormat `env:"LOG_FORMAT" envDefault:"text"`
	Health       ServerConfig     `envPrefix:"HEALTH_"`
	HealthPort   int              `env:"HEALTH_PORT" envDefault:"8080"`
	Webhook      ServerConfig     `envPrefix:"WEBHOOK_"`
	WebhookPort  int              `env:"WEBHOOK_PORT" envDefault:"8888"`
	Porkbun      PorkbunConfig    `envPrefix:"PORKBUN_"`
	DomainConfig ExternalDnsDomainConfig
}

func GetConfig() (*Config, error) {
	godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
