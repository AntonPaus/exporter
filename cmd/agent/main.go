package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
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

func pollStats(stats chan Metrics, interval time.Duration) {
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
		savedStats.PollCount += 1
		savedStats.RandomValue = gauge(rand.Float64())
		// fmt.Printf("Alloc: %.2f MB\nTotalAlloc: %.2f MB\nPollCount = %d\nRand = %.2f\n\n", savedStats.Alloc, savedStats.TotalAlloc, savedStats.PollCount, savedStats.RandomValue)
		fmt.Printf("Poll completed\n")
		stats <- savedStats
	}
}

func reportStats(stats chan Metrics, interval time.Duration, ep string) {
	// sleepDuration := time.Duration(reportInterval) * time.Second
	var receivedStats Metrics
	var valueStr string
	for {
		time.Sleep(interval)
	innerLoop:
		for {
			select {
			case r2 := <-stats:
				receivedStats = r2
				// fmt.Println("continue")
			default:
				// fmt.Printf("\n-----\nAlloc: %.2f MB\n-----\n", receivedStats.Alloc)
				statsType := reflect.TypeOf(receivedStats)
				statsValue := reflect.ValueOf(receivedStats)
				for i := 0; i < statsType.NumField(); i++ {
					field := statsType.Field(i)
					value := statsValue.Field(i)
					fieldTypeParts := strings.Split(field.Type.String(), ".")
					fieldType := fieldTypeParts[len(fieldTypeParts)-1]

					switch value.Kind() {
					case reflect.Int64:
						valueStr = fmt.Sprintf("%d", value.Int())
					case reflect.Float64:
						valueStr = fmt.Sprintf("%.2f", value.Float())
					default:
						valueStr = fmt.Sprintf("%v", value.Interface())
					}
					s := fmt.Sprintf("http://%s/update/%s/%s/%s", ep, fieldType, field.Name, valueStr)
					fmt.Printf("%s\n", s)
					req, err := http.NewRequest("POST", s, nil)
					if err != nil {
						fmt.Println("Error creating request:", err)
						return
					}
					req.Header.Set("Content-Type", "text/plain")
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						fmt.Println("Error making HTTP request:", err)
						return
					}
					defer resp.Body.Close()
					// fmt.Println("Response Status:", resp.Status)
					// fmt.Println("Response Headers:", resp.Header)
				}
				break innerLoop
			}
		}
	}
}

func main() {
	ep := flag.String("a", "localhost:8080", "server endpoint")
	ri := flag.Int("r", 10, "reportInterval")
	pi := flag.Int("p", 2, "pollInterval")
	flag.Parse()
	stats := make(chan Metrics, 6)
	go pollStats(stats, time.Duration(*pi)*time.Second)
	go reportStats(stats, time.Duration(*ri)*time.Second, *ep)
	select {}
}
