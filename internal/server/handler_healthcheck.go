package server

import "net/http"

func (h *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := h.Storage.HealthCheck()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
