package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	Address string `env:"ADDRESS"`
	// reportInterval int    `env:"REPORT_INTERVAL" envDefault:10`
	// pollInterval   int    `env:"POLL_INTERVAL" envDefault:2`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return cfg, nil
}
