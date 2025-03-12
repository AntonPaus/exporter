package memory

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/AntonPaus/exporter/internal/repository"
)

type gauge float64
type counter int64
type Memory struct {
	repository.Interfacer
	mu sync.Mutex
	G  map[string]gauge
	C  map[string]counter
	// a string
	// b int64
}

func NewMemory() *Memory {
	return &Memory{
		G: make(map[string]gauge),
		C: make(map[string]counter),
	}
}

func (m *Memory) Update(metricType string, metricName string, metricValue string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch metricType {
	case "gauge":
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("invalid gauge value")
		}
		fmt.Println("->", metricType, metricName, metricValue)
		m.G[metricName] = gauge(floatValue)
	case "counter":
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("invalid counter value")
		}
		_, ok := m.C[metricName]
		if ok {
			m.C[metricName] += counter(intValue)
		} else {
			m.C[metricName] = counter(intValue)
		}
	default:
		return errors.New("unknown metric type")
	}
	return nil
}
