package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AntonPaus/exporter/internal/storages/memory"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int
}

var m MemStorage

func mainPage(res http.ResponseWriter, req *http.Request, storage *memory.Memory) {
	// http.Error(res, "Wrong URL!", http.StatusNotFound)
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	body := fmt.Sprintf("Method: %s\r\n", req.Method)
	for k, v := range req.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "Query parameters ===============\r\n"
	for k, v := range req.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	// кодируем в JSON
	js, err := json.Marshal(storage.G)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	js2, err := json.Marshal(storage.C)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	res.Header().Set("content-type", "application/json")
	// устанавливаем код 200
	res.WriteHeader(http.StatusNotFound)
	// пишем тело ответа
	res.Write([]byte(body))
	res.Write(js)
	res.Write(js2)
}

func updateMetric(res http.ResponseWriter, req *http.Request, storage *memory.Memory) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	// check url correctness
	components := strings.Split(req.URL.Path, "/")
	if len(components) != 5 {
		http.Error(res, "Wrong URL!", http.StatusNotFound)
		return
	}
	// fmt.Println("Continue")
	// Check metric type
	switch components[2] {
	case "gauge":
		_, err := strconv.ParseFloat(components[4], 64)
		if err != nil {
			http.Error(res, "Wrong metric value!", http.StatusBadRequest)
		}
		storage.Update("gauge", components[3], components[4])
	case "counter":
		_, err := strconv.Atoi(components[4])
		if err != nil {
			http.Error(res, "Wrong metric value!", http.StatusBadRequest)
		}
		storage.Update("counter", components[3], components[4])
	default:
		http.Error(res, "Wrong metric type!", http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(req.URL.Path))
}

// type Repository interface {
// 	GetOrCreateHour(ctx context.Context, hourTime time.Time) (*Hour, error)
// 	UpdateHour(
// 		 ctx context.Context,
// 		 hourTime time.Time,
// 		 updateFn func(h *Hour) (*Hour, error),
// 	) error
// }

func main() {
	m = MemStorage{
		Gauge: map[string]float64{
			"temperature": 22.5,
			"pressure":    101.3,
		},
		Counter: map[string]int{
			"requests": 100,
			"errors":   0,
		},
	}
	storage := memory.NewMemory()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) { mainPage(w, r, storage) })
	mux.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) { updateMetric(w, r, storage) })
	// mux.HandleFunc(`/p/`, redirect)
	// mux.HandleFunc(`/a/`, http.NotFoundHandler().ServeHTTP)
	// mux.HandleFunc(`/golang/`, http.NotFoundHandler().ServeHTTP)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
