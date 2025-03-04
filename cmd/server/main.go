package main

import (
	"net/http"

	"github.com/AntonPaus/exporter/internal/handlers"
	"github.com/AntonPaus/exporter/internal/storages/memory"
)

func main() {
	storage := memory.NewMemory()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) { handlers.MainPage(w, r, storage) })
	mux.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) { handlers.UpdateMetric(w, r, storage) })
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
