package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AntonPaus/exporter/internal/storages/memory"
)

func MainPage(res http.ResponseWriter, req *http.Request, storage *memory.MemoryStorage) {
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

	// var o1 map[string]interface{}
	// var l string
	js1, err := json.Marshal(storage.G)
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
	res.WriteHeader(http.StatusOK)
	// пишем тело ответа
	res.Write(js1)
	res.Write(js2)
}

func UpdateMetric(res http.ResponseWriter, req *http.Request, storage *memory.MemoryStorage, mType string, mName string, mValue string) {
	// if req.Method != http.MethodPost {
	// 	http.Error(res, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
	// 	return
	// }
	// components := strings.Split(req.URL.Path, "/")
	// if len(components) != 5 {
	// 	http.Error(res, "Wrong URL!", http.StatusNotFound)
	// 	return
	// }
	switch mType {
	case "gauge":
		_, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(res, "Wrong metric value!", http.StatusBadRequest)
		}
	case "counter":
		_, err := strconv.Atoi(mValue)
		if err != nil {
			http.Error(res, "Wrong metric value!", http.StatusBadRequest)
		}
	default:
		http.Error(res, "Wrong metric type!", http.StatusBadRequest)
		return
	}
	storage.Update(mType, mName, mValue)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(req.URL.Path))
}

func GetMetric(w http.ResponseWriter, r *http.Request, storage *memory.MemoryStorage, mType string, mName string) {
	var valueStr string
	value, err := storage.Get(mType, mName)
	if err != nil {
		http.Error(w, "Wrong metric value!", http.StatusNotFound)
		return
	}
	switch v := value.(type) {
	case memory.Gauge:
		valueStr = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case memory.Counter:
		valueStr = strconv.FormatInt(int64(v), 10)
	default:
		http.Error(w, "Unsupported value type", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(valueStr))
}
