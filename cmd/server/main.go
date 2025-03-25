package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AntonPaus/exporter/internal/compression"
	"github.com/AntonPaus/exporter/internal/config"
	"github.com/AntonPaus/exporter/internal/handlers"
	"github.com/AntonPaus/exporter/internal/logger"
	"github.com/AntonPaus/exporter/internal/storages/memory"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Config   *config.Config
	Storage  *memory.Storage
	Router   *chi.Mux
	Handlers handlers.Handler
	// Logger          *log.Logger
}

func NewApp(cfg *config.Config) (*App, error) {
	storage, err := memory.NewStorage(cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore)
	if err != nil {
		return nil, fmt.Errorf("cannot initiate storage: %w", err)
	}
	app := &App{
		Config:  cfg,
		Storage: storage,
	}
	app.Handlers = handlers.Handler{
		Storage: app.Storage,
	}
	app.Router = chi.NewRouter()
	app.setupRoutes()
	return app, nil
}

func (a *App) setupRoutes() {
	a.Router.Use(logger.WithLogging)
	a.Router.Get("/", a.Handlers.MainPage)
	a.Router.Route("/update", func(r chi.Router) {
		r.Use(compression.WithUncompressGzip)
		r.Post("/", a.Handlers.UpdateMetricJSON)
		r.Post("/{type}/{name}/{value}", a.Handlers.UpdateMetric)
	})
	a.Router.Post("/value/", a.Handlers.GetMetricJSON)
	a.Router.Get("/value/{type}/{name}", a.Handlers.GetMetric)
}

func (a *App) Run() {
	fmt.Println("Current config:")
	fmt.Println("\tRestore values:", a.Config.Restore)
	fmt.Println("\tServer address:", a.Config.Address)
	fmt.Println("\tFile storage path:", a.Config.FileStoragePath)
	fmt.Println("\tStore interval:", a.Config.StoreInterval)
	fmt.Printf("\nStarting server on %s", a.Config.Address)
	log.Fatal(http.ListenAndServe(a.Config.Address, a.Router))
}

func main() {
	cfg, err := config.NewConfigServer()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %s", err)
	}
	defer app.Storage.Terminate()
	app.Run()
}
