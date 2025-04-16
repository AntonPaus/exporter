package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Server) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
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
		_, err := h.Storage.UpdateGauge(metrics.ID, *metrics.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case MetricTypeCounter:
		v, err := h.Storage.UpdateCounter(metrics.ID, *metrics.Delta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		*metrics.Delta = v
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
