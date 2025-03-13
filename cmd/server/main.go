package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/AntonPaus/exporter/internal/handlers"
	"github.com/AntonPaus/exporter/internal/storages/memory"
	"github.com/go-chi/chi/v5"
)

func main() {
	ep := flag.String("a", "localhost:8080", "server endpoint")
	flag.Parse()
	storage := memory.NewMemory()
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { handlers.MainPage(w, r, storage) })
	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateMetric(w, r, storage, chi.URLParam(r, "type"), chi.URLParam(r, "name"), chi.URLParam(r, "value"))
	})
	r.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMetric(w, r, storage, chi.URLParam(r, "type"), chi.URLParam(r, "name"))
	})
	log.Fatal(http.ListenAndServe(*ep, r))
}
