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
	DatabaseDSN     string `env:"DATABASE_DSN"`
	StorageType     string `env:"-"`
}

func LoadConfigServer() (Config, error) {
	cfg := Config{}
	if err := env.Parse(cfg); err != nil {
		return Config{}, fmt.Errorf("config error: %w", err)
	}
	address := new(string)
	fileStoragePath := new(string)
	restore := new(bool)
	storeInterval := new(uint)
	databaseDSN := new(string)
	flag.StringVar(address, "a", "localhost:8080", "server endpoint")
	flag.UintVar(storeInterval, "i", 300, "Store interval")
	flag.StringVar(fileStoragePath, "f", "./storage.json", "f")
	flag.BoolVar(restore, "r", false, "restore config")
	flag.StringVar(databaseDSN, "d", "localhost", "database endpoint")
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
	if cfg.DatabaseDSN == "" {
		cfg.DatabaseDSN = *databaseDSN
	}
	switch {
	case cfg.DatabaseDSN != "":
		cfg.StorageType = "database"
	case cfg.FileStoragePath != "":
		cfg.StorageType = "file"
	default:
		cfg.StorageType = "memory"
	}
	fmt.Println("Current config:")
	fmt.Println("\tRestore values:", cfg.Restore)
	fmt.Println("\tServer address:", cfg.Address)
	fmt.Println("\tFile storage path:", cfg.FileStoragePath)
	fmt.Println("\tStore interval:", cfg.StoreInterval)
	fmt.Println("\tDB address:", cfg.DatabaseDSN)
	fmt.Println("\tStorage type:", cfg.StorageType)
	return cfg, nil
}
