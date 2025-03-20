package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AntonPaus/exporter/internal/config"
	"github.com/AntonPaus/exporter/internal/handlers"
	"github.com/AntonPaus/exporter/internal/storages/memory"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	Storage memory.MemoryStorage
}

var sugar zap.SugaredLogger

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	responseData        *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

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
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
	sugar.Infow(
		"Starting server",
		"addr", *address,
	)

	h := handlers.Handler{
		Storage: storage,
	}

	r := chi.NewRouter()
	r.Use(WithLogging)
	// r.Get("/", h.MainPage)
	r.Post("/update/", h.UpdateMetricJSON)
	r.Post("/value/", h.GetMetricJSON)
	r.Get("/value/{type}/{name}", h.GetMetric)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	log.Fatal(http.ListenAndServe(*address, r))
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter
		duration := time.Since(start)
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"data", string(bodyBytes),
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn)
}
