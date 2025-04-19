package app

import (
	"fmt"

	"github.com/AntonPaus/exporter/internal/config"
	"github.com/AntonPaus/exporter/internal/server"
	"github.com/AntonPaus/exporter/internal/storage"
	"github.com/AntonPaus/exporter/internal/storage/memory"
)

type App struct {
	server *server.Server
	ep     string
}

func NewApp(cfg *config.Config) (*App, error) {
	var storage storage.Storage
	var err error
	switch cfg.StorageType {
	case "database":
		storage, err = memory.NewStorage(cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore)
	case "file":
		// storage = file.NewFileStorage(cfg.FilePath)
		storage, err = memory.NewStorage(cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore)
	case "memory":
		// storage = database.NewDatabaseStorage(cfg.DatabaseDSN)
		storage, err = memory.NewStorage(cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.StorageType)
	}
	if err != nil {
		return nil, err
	}
	srv := server.NewServer(storage)
	return &App{server: srv, ep: cfg.Address}, nil
}

func (a *App) Run() error {
	return a.server.Start(a.ep)
}
