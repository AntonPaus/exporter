package memory

import (
	"errors"
	"strconv"
	"sync"

	"github.com/AntonPaus/exporter/internal/repository"
)

type Gauge float64
type Counter int64
type Memory struct {
	repository.Interfacer
	mu sync.Mutex
	G  map[string]Gauge
	C  map[string]Counter
	// a string
	// b int64
}

func NewMemory() *Memory {
	return &Memory{
		G: make(map[string]Gauge),
		C: make(map[string]Counter),
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

		m.G[metricName] = Gauge(floatValue)
	case "counter":
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("invalid counter value")
		}
		_, ok := m.C[metricName]
		if ok {
			m.C[metricName] += Counter(intValue)
		} else {
			m.C[metricName] = Counter(intValue)
		}
	default:
		return errors.New("unknown metric type")
	}
	return nil
}

func (m *Memory) Get(metricType string, metricName string) (any, error) {
	switch metricType {
	case "gauge":
		v, ok := m.G[metricName]
		if ok {
			return v, nil
		}
		return nil, errors.New("wrong gauge metric name")
	case "counter":
		v, ok := m.C[metricName]
		if ok {
			return v, nil
		}
		return nil, errors.New("wrong counter metric name")
	default:
		return nil, errors.New("unknown metric type")
	}
}
