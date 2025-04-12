package main

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/AntonPaus/exporter/internal/app"
	"github.com/AntonPaus/exporter/internal/config"
)

func main() {
	cfg, err := config.LoadConfigServer()
	fmt.Println(cfg)
	if err != nil {
		log.Fatalf("Failed to parse arguments. Full error: %v", err)
	}
	fmt.Println("Current ")
	// storageType := flag.String("storage", "memory", "storage type: memory, file, or database")
	// flag.Parse()
	fmt.Println(cfg)
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application stopped with error: %v", err)
	}
}
