package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int
}

var m MemStorage

func mainPage(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Wrong URL!", http.StatusNotFound)
	// err := req.ParseForm()
	// if err != nil {
	// 	panic(err)
	// }
	// body := fmt.Sprintf("Method: %s\r\n", req.Method)
	// for k, v := range req.Header {
	// 	body += fmt.Sprintf("%s: %v\r\n", k, v)
	// }
	// body += "Query parameters ===============\r\n"
	// for k, v := range req.Form {
	// 	body += fmt.Sprintf("%s: %v\r\n", k, v)
	// }
	// кодируем в JSON
	// resp, err := json.Marshal(m)
	// if err != nil {
	// 	http.Error(res, err.Error(), 500)
	// 	return
	// }
	// res.Header().Set("content-type", "application/json")
	// устанавливаем код 200
	// res.WriteHeader(http.StatusNotFound)
	// пишем тело ответа
	// res.Write([]byte(body))
	// res.Write(resp)
}

func updateMetric(res http.ResponseWriter, req *http.Request) {
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
	if components[2] == "gauge" {
		number, err := strconv.ParseFloat(components[4], 64)
		if err != nil {
			http.Error(res, "Wrong metric value!", http.StatusBadRequest)
		}
		m.Gauge[components[3]] = number
	} else if components[2] == "counter" {
		number, err := strconv.Atoi(components[4])
		if err != nil {
			http.Error(res, "Wrong metric value!", http.StatusBadRequest)
		}
		_, ok := m.Counter[components[3]]
		if ok {
			m.Counter[components[3]] += number
		} else {
			m.Counter[components[3]] = number
		}
	} else {
		http.Error(res, "Wrong metric type!", http.StatusBadRequest)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(req.URL.Path))
}

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
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	mux.HandleFunc(`/update/`, updateMetric)
	// mux.HandleFunc(`/p/`, redirect)
	// mux.HandleFunc(`/a/`, http.NotFoundHandler().ServeHTTP)
	// mux.HandleFunc(`/golang/`, http.NotFoundHandler().ServeHTTP)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
