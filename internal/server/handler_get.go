package server

import (
	"net/http"
	"strconv"

	"github.com/AntonPaus/exporter/internal/server/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetMetric(w http.ResponseWriter, r *http.Request) {
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
		resp, err = middleware.CompressGzip([]byte(valueStr))
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
