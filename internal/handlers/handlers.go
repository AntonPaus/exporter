package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/AntonPaus/exporter/internal/compression"
	"github.com/go-chi/chi/v5"
)

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type Storage interface {
	Get(mType string, mName string) (any, error)
	Update(mType string, mName string, mValue any) (any, error)
	Terminate()
}

type Handler struct {
	Storage Storage
	DB      *sql.DB
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	err := h.DB.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	resp := []byte("Fake body")
	if r.Header.Get("Accept-Encoding") == "gzip" {
		resp, err = compression.CompressGzip([]byte(resp))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics
	var ok bool
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch metrics.MType {
	case MetricTypeGauge:
		v, err := h.Storage.Update(metrics.ID, metrics.MType, *metrics.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if *metrics.Value, ok = v.(float64); !ok {
			http.Error(w, "value is not gauge type", http.StatusBadRequest)
			return
		}
	case MetricTypeCounter:
		v, err := h.Storage.Update(metrics.ID, metrics.MType, *metrics.Delta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if *metrics.Delta, ok = v.(int64); !ok {
			http.Error(w, "value is not counter type", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Unsupported value type", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	value, err := h.Storage.Get(metrics.ID, metrics.MType)
	if err != nil {
		http.Error(w, "Wrong metric value!", http.StatusNotFound)
		return
	}
	switch metrics.MType {
	case MetricTypeGauge:
		if value, ok := value.(float64); ok {
			metrics.Value = &value
		} else {
			http.Error(w, "Unsupported value type", http.StatusInternalServerError)
			return
		}
	case MetricTypeCounter:
		if value, ok := value.(int64); ok {
			metrics.Delta = &value
		} else {
			http.Error(w, "Unsupported value type", http.StatusInternalServerError)
			return
		}
		*metrics.Delta = value.(int64)
	default:
		http.Error(w, "Unsupported value type", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Header.Get("Accept-Encoding") == "gzip" {
		resp, err = compression.CompressGzip(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var mType, mName, mValue string
	mType = chi.URLParam(r, "type")
	mName = chi.URLParam(r, "name")
	mValue = chi.URLParam(r, "value")
	fmt.Println(mValue)
	switch mType {
	case MetricTypeGauge:
		g, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, "Wrong metric value!", http.StatusBadRequest)
			return
		}
		h.Storage.Update(mName, mType, g)
	case MetricTypeCounter:
		c, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value!", http.StatusBadRequest)
			return
		}
		h.Storage.Update(mName, mType, c)
	default:
		http.Error(w, "Wrong metric type!", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.Path))
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	var mType, mName, valueStr string

	mType = chi.URLParam(r, "type")
	mName = chi.URLParam(r, "name")
	value, err := h.Storage.Get(mName, mType)
	if err != nil {
		http.Error(w, "Wrong metric value!", http.StatusNotFound)
		return
	}
	switch mType {
	case MetricTypeGauge:
		if value, ok := value.(float64); ok {
			valueStr = strconv.FormatFloat(float64(value), 'f', -1, 64)
		} else {
			http.Error(w, "Unsupported value type", http.StatusInternalServerError)
			return
		}
	case MetricTypeCounter:
		if value, ok := value.(int64); ok {
			valueStr = strconv.FormatInt(int64(value), 10)
		} else {
			http.Error(w, "Unsupported value type", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsupported value type", http.StatusInternalServerError)
		return
	}
	resp := []byte(valueStr)
	if r.Header.Get("Accept-Encoding") == "gzip" {
		resp, err = compression.CompressGzip([]byte(valueStr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
