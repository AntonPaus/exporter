package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/AntonPaus/exporter/internal/server/middleware"
)

type jsonMetrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
type gauge float64
type counter int64

type Metrics struct {
	values struct {
		Alloc         gauge
		BuckHashSys   gauge
		Frees         gauge
		GCCPUFraction gauge
		GCSys         gauge
		HeapAlloc     gauge
		HeapIdle      gauge
		HeapInuse     gauge
		HeapObjects   gauge
		HeapReleased  gauge
		HeapSys       gauge
		LastGC        gauge
		Lookups       gauge
		MCacheInuse   gauge
		MCacheSys     gauge
		MSpanInuse    gauge
		MSpanSys      gauge
		Mallocs       gauge
		NextGC        gauge
		NumForcedGC   gauge
		NumGC         gauge
		OtherSys      gauge
		PauseTotalNs  gauge
		StackInuse    gauge
		StackSys      gauge
		Sys           gauge
		TotalAlloc    gauge
		RandomValue   gauge
		PollCount     counter
	}
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Poll(interval time.Duration) {
	var mem runtime.MemStats
	for {
		time.Sleep(interval)
		runtime.ReadMemStats(&mem)
		m.values.Alloc = gauge(mem.Alloc)
		m.values.BuckHashSys = gauge(mem.BuckHashSys)
		m.values.Frees = gauge(mem.Frees)
		m.values.GCCPUFraction = gauge(mem.GCCPUFraction)
		m.values.GCSys = gauge(mem.GCSys)
		m.values.HeapAlloc = gauge(mem.HeapAlloc)
		m.values.HeapIdle = gauge(mem.HeapIdle)
		m.values.HeapInuse = gauge(mem.HeapInuse)
		m.values.HeapObjects = gauge(mem.HeapObjects)
		m.values.HeapReleased = gauge(mem.HeapReleased)
		m.values.HeapSys = gauge(mem.HeapSys)
		m.values.LastGC = gauge(mem.LastGC)
		m.values.Lookups = gauge(mem.Lookups)
		m.values.MCacheInuse = gauge(mem.MCacheInuse)
		m.values.MCacheSys = gauge(mem.MCacheSys)
		m.values.MSpanInuse = gauge(mem.MSpanInuse)
		m.values.MSpanSys = gauge(mem.MSpanSys)
		m.values.Mallocs = gauge(mem.Mallocs)
		m.values.NextGC = gauge(mem.NextGC)
		m.values.NumForcedGC = gauge(mem.NumForcedGC)
		m.values.NumGC = gauge(mem.NumGC)
		m.values.OtherSys = gauge(mem.OtherSys)
		m.values.PauseTotalNs = gauge(mem.PauseTotalNs)
		m.values.StackInuse = gauge(mem.StackInuse)
		m.values.StackSys = gauge(mem.StackSys)
		m.values.Sys = gauge(mem.Sys)
		m.values.TotalAlloc = gauge(mem.TotalAlloc)
		m.values.PollCount = 1
		m.values.RandomValue = gauge(rand.Float64())
		fmt.Printf("Poll completed\n")
	}
}

func (m *Metrics) Report(interval time.Duration, ep string) {
	var c int64
	var g float64
	for {
		errFound := false
		time.Sleep(interval)
		statsType := reflect.TypeOf(m.values)
		statsValue := reflect.ValueOf(m.values)
		for i := range statsType.NumField() {
			var h jsonMetrics
			field := statsType.Field(i)
			value := statsValue.Field(i)
			fieldTypeParts := strings.Split(field.Type.String(), ".")
			fieldType := fieldTypeParts[len(fieldTypeParts)-1]
			h.ID, h.MType = field.Name, fieldType
			switch value.Kind() {
			case reflect.Int64:
				c = int64(value.Int())
				h.Delta = &c
			case reflect.Float64:
				g = float64(value.Float())
				h.Value = &g
			default:
				fmt.Printf("Value type error\nSkipping...\n")
			}
			jsonData, err := json.Marshal(h)
			if err != nil {
				fmt.Printf("JSON Marshaling error: %v\nSkipping...\n", err)
			}
			compressedData, err := middleware.CompressGzip(jsonData)
			if err != nil {
				fmt.Printf("Compression error: %v\nSkipping...\n", err)
			}
			if err := sendMetric(compressedData, ep); err != nil {
				fmt.Println("Error sending HTTP request:", err)
				errFound = true
				break
			}
		}
		if !errFound {
			fmt.Println("Report completed")
		}
	}
}

func sendMetric(compressedData []byte, ep string) error {
	s := fmt.Sprintf("http://%s/update/", ep)
	req, err := http.NewRequest("POST", s, bytes.NewBuffer(compressedData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	return nil
}
