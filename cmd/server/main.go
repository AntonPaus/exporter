package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/AntonPaus/exporter/internal/compression"
	"github.com/AntonPaus/exporter/internal/config"
	"github.com/AntonPaus/exporter/internal/database"
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
	db       *sql.DB
	// Logger          *log.Logger
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	migrations, err := database.LoadMigrations("db/migrations")
	if err != nil {
		return nil, err
	}
	err = database.ApplyMigrations(db, migrations)
	if err != nil {
		return nil, err
	}
	storage, err := memory.NewStorage(cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore)
	if err != nil {
		return nil, err
	}
	app := &App{
		Config:  cfg,
		Storage: storage,
		db:      db,
	}
	app.Handlers = handlers.Handler{
		Storage: app.Storage,
		DB:      app.db,
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
	a.Router.Get("/ping", a.Handlers.PingDB)
}

func (a *App) Run() {
	fmt.Println("Current config:")
	fmt.Println("\tRestore values:", a.Config.Restore)
	fmt.Println("\tServer address:", a.Config.Address)
	fmt.Println("\tFile storage path:", a.Config.FileStoragePath)
	fmt.Println("\tStore interval:", a.Config.StoreInterval)
	fmt.Println("\tDB address:", a.Config.DatabaseDSN)
	fmt.Printf("\nStarting server on %s ", a.Config.Address)
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
	defer app.db.Close()
	defer app.Storage.Terminate()
	app.Run()
}
