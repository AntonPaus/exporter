package server

import (
	"net/http"

	"github.com/AntonPaus/exporter/internal/logger"
	"github.com/AntonPaus/exporter/internal/server/middleware"
	"github.com/AntonPaus/exporter/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
	Storage storage.Storage
	logger  *zap.SugaredLogger
}

func NewServer(storage storage.Storage) *Server {
	sugar := logger.GetLogger()
	s := &Server{
		router:  chi.NewRouter(),
		Storage: storage,
		logger:  sugar,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.Use(logger.WithLoggingNew)
	s.router.Get("/", s.MainPage)
	s.router.Route("/update", func(r chi.Router) {
		r.Use(middleware.WithUncompressGzip)
		r.Post("/", s.UpdateMetricJSON)
		r.Post("/{type}/{name}/{value}", s.UpdateMetric)
	})
	s.router.Post("/value/", s.GetMetricJSON)
	s.router.Get("/value/{type}/{name}", s.GetMetric)
	s.router.Get("/ping", s.HealthCheck)
}

func (s *Server) Start(ep string) error {
	s.logger.Infow("Starting server on ", "address", ep)
	return http.ListenAndServe(ep, s.router)
}
