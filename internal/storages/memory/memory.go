package memory

import (
	"errors"
	"strconv"
	"sync"

	"github.com/AntonPaus/exporter/internal/repository"
)

type Memory struct {
	repository.Interfacer
	mu sync.Mutex
	g  map[string]float64
	c  map[string]int64
	// a string
	// b int64
}

func NewMemory() *Memory {
	return &Memory{
		g: make(map[string]float64),
		c: make(map[string]int64),
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
		m.g[metricName] = floatValue
	case "counter":
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("invalid counter value")
		}
		m.c[metricName] = intValue
	default:
		return errors.New("unknown metric type")
	}
	return nil
}
