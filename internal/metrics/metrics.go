package metrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/AntonPaus/exporter/internal/compression"
	"github.com/AntonPaus/exporter/internal/handlers"
)

type gauge float64
type counter int64

type Metrics struct {
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

func PollStats(stats chan Metrics, interval time.Duration) {
	var mem runtime.MemStats
	var savedStats Metrics
	savedStats.PollCount = 0
	for {
		time.Sleep(interval)
		runtime.ReadMemStats(&mem)
		savedStats.Alloc = gauge(mem.Alloc)
		savedStats.BuckHashSys = gauge(mem.BuckHashSys)
		savedStats.Frees = gauge(mem.Frees)
		savedStats.GCCPUFraction = gauge(mem.GCCPUFraction)
		savedStats.GCSys = gauge(mem.GCSys)
		savedStats.HeapAlloc = gauge(mem.HeapAlloc)
		savedStats.HeapIdle = gauge(mem.HeapIdle)
		savedStats.HeapInuse = gauge(mem.HeapInuse)
		savedStats.HeapObjects = gauge(mem.HeapObjects)
		savedStats.HeapReleased = gauge(mem.HeapReleased)
		savedStats.HeapSys = gauge(mem.HeapSys)
		savedStats.LastGC = gauge(mem.LastGC)
		savedStats.Lookups = gauge(mem.Lookups)
		savedStats.MCacheInuse = gauge(mem.MCacheInuse)
		savedStats.MCacheSys = gauge(mem.MCacheSys)
		savedStats.MSpanInuse = gauge(mem.MSpanInuse)
		savedStats.MSpanSys = gauge(mem.MSpanSys)
		savedStats.Mallocs = gauge(mem.Mallocs)
		savedStats.NextGC = gauge(mem.NextGC)
		savedStats.NumForcedGC = gauge(mem.NumForcedGC)
		savedStats.NumGC = gauge(mem.NumGC)
		savedStats.OtherSys = gauge(mem.OtherSys)
		savedStats.PauseTotalNs = gauge(mem.PauseTotalNs)
		savedStats.StackInuse = gauge(mem.StackInuse)
		savedStats.StackSys = gauge(mem.StackSys)
		savedStats.Sys = gauge(mem.Sys)
		savedStats.TotalAlloc = gauge(mem.TotalAlloc)
		savedStats.PollCount = 1
		savedStats.RandomValue = gauge(rand.Float64())
		fmt.Printf("Poll completed\n")
		stats <- savedStats
	}
}

func ReportStats(stats chan Metrics, interval time.Duration, ep string) error {
	var receivedStats Metrics
	var c int64
	var g float64
	// var valueStr string
	for {
		time.Sleep(interval)
	innerLoop:
		for {
			errFound := false
			select {
			case r2 := <-stats:
				receivedStats = r2
			default:
				statsType := reflect.TypeOf(receivedStats)
				statsValue := reflect.ValueOf(receivedStats)
				for i := 0; i < statsType.NumField(); i++ {
					var m handlers.Metrics
					field := statsType.Field(i)
					value := statsValue.Field(i)
					fieldTypeParts := strings.Split(field.Type.String(), ".")
					fieldType := fieldTypeParts[len(fieldTypeParts)-1]
					m.ID, m.MType = field.Name, fieldType
					switch value.Kind() {
					case reflect.Int64:
						c = int64(value.Int())
						m.Delta = &c
						// valueStr = fmt.Sprintf("%d", value.Int())
					case reflect.Float64:
						g = float64(value.Float())
						m.Value = &g
					default:
						return errors.New("unsupported type")
					}
					// fmt.Println(m)
					jsonData, err := json.Marshal(m)
					if err != nil {
						return err
					}
					compressedData, err := compression.CompressGzip(jsonData)
					if err != nil {
						return err
					}
					s := fmt.Sprintf("http://%s/update/", ep)
					// fmt.Printf("%s\n", s)

					req, err := http.NewRequest("POST", s, bytes.NewBuffer(compressedData))
					if err != nil {
						fmt.Println("Error creating request:", err)
						return err
					}
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Content-Encoding", "gzip")
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						fmt.Println("Error making HTTP request:", err)
						errFound = true
						break
					}
					defer resp.Body.Close()
				}
				if !errFound {
					fmt.Println("Sending completed")
				}
				break innerLoop
			}
		}
	}
}
