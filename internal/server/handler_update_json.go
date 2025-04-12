package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (h *Handlers) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
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
		vStr := fmt.Sprintf("%v", *metrics.Value)
		v, err := h.Storage.Update(metrics.ID, metrics.MType, vStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if *metrics.Value, ok = v.(float64); !ok {
			http.Error(w, "value is not gauge type", http.StatusBadRequest)
			return
		}
	case MetricTypeCounter:
		vStr := fmt.Sprintf("%v", *metrics.Value)
		v, err := h.Storage.Update(metrics.ID, metrics.MType, vStr)
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
