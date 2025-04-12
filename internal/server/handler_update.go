package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var mType, mName, mValue string
	mType = chi.URLParam(r, "type")
	mName = chi.URLParam(r, "name")
	mValue = chi.URLParam(r, "value")
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
		_, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, "Wrong metric value!", http.StatusBadRequest)
			return
		}
		h.Storage.Update(mName, mType, mValue)

	case MetricTypeCounter:
		_, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value!", http.StatusBadRequest)
			return
		}
		h.Storage.Update(mName, mType, mValue)
	default:
		http.Error(w, "Wrong metric type!", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.Path))
}
