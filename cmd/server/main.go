package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/AntonPaus/exporter/internal/config"
	"github.com/AntonPaus/exporter/internal/handlers"
	"github.com/AntonPaus/exporter/internal/storages/memory"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	address := new(string)
	flag.StringVar(address, "a", "localhost:8080", "server endpoint")
	flag.Parse()
	if cfg.Address != "" {
		*address = cfg.Address
	}
	storage := memory.NewMemoryStorage()
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { handlers.MainPage(w, r, storage) })
	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateMetric(w, r, storage, chi.URLParam(r, "type"), chi.URLParam(r, "name"), chi.URLParam(r, "value"))
	})
	r.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMetric(w, r, storage, chi.URLParam(r, "type"), chi.URLParam(r, "name"))
	})
	log.Fatal(http.ListenAndServe(*address, r))
}
