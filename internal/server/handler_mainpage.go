package server

import (
	"net/http"

	"github.com/AntonPaus/exporter/internal/server/middleware"
)

func (h *Server) MainPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	resp := []byte("Fake body")
	if r.Header.Get("Accept-Encoding") == "gzip" {
		resp, err = middleware.CompressGzip([]byte(resp))
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
