package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"
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
	r.Get("/", h.MainPage)
	r.Route("/update", func(r chi.Router) {
		r.Use(WithUncompressGzip)
		r.Post("/", h.UpdateMetricJSON)
		r.Post("/{type}/{name}/{value}", h.UpdateMetric)
	})
	// r.Post("/update/", h.UpdateMetricJSON)
	// r.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	r.Post("/value/", h.GetMetricJSON)
	r.Get("/value/{type}/{name}", h.GetMetric)
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
			"compression", r.Header.Get("Content-Encoding"),
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"data", string(bodyBytes),
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn)
}

func WithUncompressGzip(h http.Handler) http.Handler {
	uncompressFn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// закрытие gzip-читателя опционально, т.к. все данные уже прочитаны и
			// текущая реализация не требует закрытия, тем не менее лучше это делать -
			// некоторые реализации могут рассчитывать на закрытие читателя
			// gz.Close() не вызывает закрытия r.Body - это будет сделано позже, http-сервером
			defer gz.Close()

			// при чтении вернётся распакованный слайс байт
			decompressedBody, err := io.ReadAll(gz)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(decompressedBody))
			r.ContentLength = int64(len(decompressedBody))
			r.Header.Del("Content-Encoding")
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(uncompressFn)
}
