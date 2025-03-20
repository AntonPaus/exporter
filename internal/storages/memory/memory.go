package memory

import (
	"errors"
	"sync"
)

type gauge float64
type counter int64

type MemoryStorage struct {
	mu sync.Mutex
	g  map[string]gauge
	c  map[string]counter
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		g: make(map[string]gauge),
		c: make(map[string]counter),
	}
}

func (m *MemoryStorage) Update(mName string, mType string, mValue any) (any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch mType {
	case "gauge":
		m.g[mName] = gauge(mValue.(float64))
		return float64(m.g[mName]), nil
	case "counter":
		m.c[mName] += counter(mValue.(int64))
		return int64(m.c[mName]), nil
	default:
		return nil, errors.New("something went wrong")
	}
}

func (m *MemoryStorage) Get(mName string, mType string) (any, error) {
	switch mType {
	case "gauge":
		v, ok := m.g[mName]
		if ok {
			return float64(v), nil
		}
		return nil, errors.New("no metric found")
	case "counter":
		v, ok := m.c[mName]
		if ok {
			return int64(v), nil
		}
		return nil, errors.New("no metric found")
	default:
		return nil, errors.New("unknown metric type")
	}
}
