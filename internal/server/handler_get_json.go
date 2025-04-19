package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AntonPaus/exporter/internal/server/middleware"
)

func (s *Server) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
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
	switch metrics.MType {
	case MetricTypeGauge:
		v, err := s.Storage.GetGauge(metrics.ID)
		if err != nil {
			http.Error(w, "Wrong metric!", http.StatusInternalServerError)
			return
		}
		*metrics.Value = float64(v)
	case MetricTypeCounter:
		v, err := s.Storage.GetCounter(metrics.ID)
		if err != nil {
			http.Error(w, "Wrong metric!", http.StatusInternalServerError)
			return
		}
		*metrics.Delta = int64(v)
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
		resp, err = middleware.CompressGzip(resp)
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
