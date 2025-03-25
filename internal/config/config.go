package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	StoreInterval   uint   `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
	// reportInterval int    `env:"REPORT_INTERVAL" envDefault:10`
	// pollInterval   int    `env:"POLL_INTERVAL" envDefault:2`
}

func NewConfigServer() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	address := new(string)
	fileStoragePath := new(string)
	restore := new(bool)
	storeInterval := new(uint)
	flag.StringVar(address, "a", "localhost:8080", "server endpoint")
	flag.UintVar(storeInterval, "i", 300, "Store interval")
	flag.StringVar(fileStoragePath, "f", "./storage", "f")
	flag.BoolVar(restore, "r", false, "restore config")
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *address
	}
	if cfg.StoreInterval == 0 {
		cfg.StoreInterval = *storeInterval
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *fileStoragePath
	}
	if !cfg.Restore {
		cfg.Restore = *restore
	}
	return cfg, nil
}
