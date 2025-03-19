package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AntonPaus/exporter/internal/compression"
	"github.com/go-chi/chi/v5"
)

type Repository interface {
	Get(mType string, mName string) (any, error)
	Update(mType string, mName string, mValue any) (any, error)
}

type Handler struct {
	Storage Repository
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	// http.Error(w, "Wrong URL!", http.StatusNotFound)
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	body := fmt.Sprintf("Method: %s\r\n", r.Method)
	for k, v := range r.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "Query parameters ===============\r\n"
	for k, v := range r.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	// // кодируем в JSON

	// // var o1 map[string]interface{}
	// // var l string
	// js1, err := json.Marshal(h.Storage.g)
	// if err != nil {
	// 	http.Error(res, err.Error(), 500)
	// 	return
	// }
	// js2, err := json.Marshal(h.Storage.)
	// if err != nil {
	// 	http.Error(res, err.Error(), 500)
	// 	return
	// }
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
	var buf bytes.Buffer
	var metrics Metrics
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch metrics.MType {
	case "gauge":
		v, err := h.Storage.Update(metrics.ID, metrics.MType, *metrics.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		*metrics.Value = v.(float64)
	case "counter":

		v, err := h.Storage.Update(metrics.ID, metrics.MType, *metrics.Delta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		*metrics.Delta = v.(int64)
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
	var buf bytes.Buffer
	var metrics Metrics
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := h.Storage.Get(metrics.ID, metrics.MType)
	if err != nil {
		http.Error(w, "Wrong metric value!", http.StatusNotFound)
		return
	}
	switch metrics.MType {
	case "gauge":
		if value, ok := value.(float64); ok {
			metrics.Value = &value
		} else {
			http.Error(w, "Unsupported value type", http.StatusInternalServerError)
			return
		}
	case "counter":
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
	case "gauge":
		g, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, "Wrong metric value!", http.StatusBadRequest)
			return
		}
		h.Storage.Update(mName, mType, g)
	case "counter":
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
	fmt.Println(value)
	if err != nil {
		http.Error(w, "Wrong metric value!", http.StatusNotFound)
		return
	}
	switch mType {
	case "gauge":
		if value, ok := value.(float64); ok {
			valueStr = strconv.FormatFloat(float64(value), 'f', -1, 64)
		} else {
			http.Error(w, "Unsupported value type", http.StatusInternalServerError)
			return
		}
	case "counter":
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
