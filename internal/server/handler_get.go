package server

import (
	"net/http"
	"strconv"

	"github.com/AntonPaus/exporter/internal/server/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Server) GetMetric(w http.ResponseWriter, r *http.Request) {
	var mType, mName, valueStr string
	mType = chi.URLParam(r, "type")
	mName = chi.URLParam(r, "name")
	switch mType {
	case MetricTypeGauge:
		value, err := h.Storage.GetGauge(mName)
		if err != nil {
			http.Error(w, "Wrong metric!", http.StatusNotFound)
			return
		}
		valueStr = strconv.FormatFloat(value, 'f', -1, 64)
	case MetricTypeCounter:
		value, err := h.Storage.GetCounter(mName)
		if err != nil {
			http.Error(w, "Wrong metric!", http.StatusNotFound)
			return
		}
		valueStr = strconv.FormatInt(value, 10)
	default:
		http.Error(w, "Unsupported value type", http.StatusInternalServerError)
		return
	}
	resp := []byte(valueStr)
	var err error
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
