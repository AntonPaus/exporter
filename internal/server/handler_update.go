package server

import (
	"net/http"
	"strconv"

	"github.com/AntonPaus/exporter/internal/storage"
	"github.com/go-chi/chi/v5"
)

func (s *Server) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var mType, mName, mValue string
	mType = chi.URLParam(r, "type")
	mName = chi.URLParam(r, "name")
	mValue = chi.URLParam(r, "value")
	s.logger.Infow("attempting to update metric",
		"type", mType,
		"name", mName,
		"value", mValue,
	)
	switch {
	case mType == "":
		http.Error(w, "Wrong metric type!", http.StatusBadRequest)
		return
	case mName == "":
		http.Error(w, "Wrong metric name!", http.StatusBadRequest)
		return
	case mValue == "":
		http.Error(w, "Wrong metric value!", http.StatusBadRequest)
		return
	default:
	}
	switch mType {
	case MetricTypeGauge:
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, "Wrong metric value!", http.StatusBadRequest)
			return
		}
		v, err := s.Storage.UpdateGauge(mName, storage.Gauge(val))
		if err != nil {
			http.Error(w, "Metric update didn't succed", http.StatusBadRequest)
			return
		}
		s.logger.Infow("updated value",
			"value", v,
		)
	case MetricTypeCounter:
		val, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value!", http.StatusBadRequest)
			return
		}
		v, err := s.Storage.UpdateCounter(mName, storage.Counter(val))
		if err != nil {
			http.Error(w, "Metric update didn't succed", http.StatusBadRequest)
			return
		}
		s.logger.Infow("updated value",
			"value", v,
		)
	default:
		http.Error(w, "Wrong metric type!", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.Path))
}
