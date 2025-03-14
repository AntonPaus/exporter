package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AntonPaus/exporter/internal/handlers"
	"github.com/AntonPaus/exporter/internal/storages/memory"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

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
	address := new(string)
	flag.StringVar(address, "a", "localhost:8080", "server endpoint")
	flag.Parse()
	osEP := os.Getenv("ADDRESS")
	if osEP != "" {
		*address = osEP
	}
	storage := memory.NewMemory()
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
	r := chi.NewRouter()
	r.Use(WithLogging)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { handlers.MainPage(w, r, storage) })
	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateMetric(w, r, storage, chi.URLParam(r, "type"), chi.URLParam(r, "name"), chi.URLParam(r, "value"))
	})
	r.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMetric(w, r, storage, chi.URLParam(r, "type"), chi.URLParam(r, "name"))
	})
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
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn)
}
