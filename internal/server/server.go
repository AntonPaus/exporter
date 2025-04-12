package server

import (
	"log"
	"net/http"

	"github.com/AntonPaus/exporter/internal/logger"
	"github.com/AntonPaus/exporter/internal/server/middleware"
	"github.com/AntonPaus/exporter/internal/storage"
	"github.com/go-chi/chi/v5"
)

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Server struct {
	router  *chi.Mux
	storage storage.Storage
}

type Handlers struct {
	Storage storage.Storage
}

func NewServer(storage storage.Storage) *Server {
	s := &Server{
		router:  chi.NewRouter(),
		storage: storage,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	handlers := &Handlers{Storage: s.storage}

	s.router.Use(logger.WithLogging)
	s.router.Get("/", handlers.MainPage)
	s.router.Route("/update", func(r chi.Router) {
		r.Use(middleware.WithUncompressGzip)
		r.Post("/", handlers.UpdateMetricJSON)
		r.Post("/{type}/{name}/{value}", handlers.UpdateMetric)
	})
	s.router.Post("/value/", handlers.GetMetricJSON)
	s.router.Get("/value/{type}/{name}", handlers.GetMetric)
	s.router.Get("/ping", handlers.HealthCheck)
}

func (s *Server) Start() error {
	log.Println("Starting server on :8080")
	return http.ListenAndServe(":8080", s.router)
}
