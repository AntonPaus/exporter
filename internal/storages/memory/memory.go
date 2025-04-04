package memory

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type Settings struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}
type gauge float64
type counter int64

type Storage struct {
	mu           sync.Mutex
	g            map[string]gauge
	c            map[string]counter
	dumpInterval uint
	dumpFile     *os.File
}

func NewStorage(dumpInterval uint, dumpLocation string, restoreFromFile bool) (*Storage, error) {
	m := &Storage{
		g:            make(map[string]gauge),
		c:            make(map[string]counter),
		dumpInterval: dumpInterval,
		dumpFile:     nil,
	}
	if restoreFromFile {
		err := m.readFromFile(dumpLocation)
		if err != nil {
			fmt.Println("No storage file found. Continue")
			// return nil, err
		}
	}
	dumpFile, err := os.OpenFile(dumpLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	m.dumpFile = dumpFile
	if m.dumpInterval > 0 {
		go m.tickerDump()
	}
	return m, nil
}

func (m *Storage) Terminate() {
	m.dumpFile.Close()
}

func (m *Storage) Update(mName string, mType string, mValue any) (any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch mType {
	case MetricTypeGauge:
		m.g[mName] = gauge(mValue.(float64))
		if m.dumpInterval == 0 {
			m.dump()
		}
		return float64(m.g[mName]), nil
	case MetricTypeCounter:
		m.c[mName] += counter(mValue.(int64))
		if m.dumpInterval == 0 {
			m.dump()
		}
		return int64(m.c[mName]), nil
	default:
		return nil, errors.New("something went wrong")
	}
}

func (m *Storage) Get(mName string, mType string) (any, error) {
	switch mType {
	case MetricTypeGauge:
		v, ok := m.g[mName]
		if ok {
			return float64(v), nil
		}
		return nil, errors.New("no metric found")
	case MetricTypeCounter:
		v, ok := m.c[mName]
		if ok {
			return int64(v), nil
		}
		return nil, errors.New("no metric found")
	default:
		return nil, errors.New("unknown metric type")
	}
}

func (m *Storage) dump() {
	var buf bytes.Buffer
	d1, err := json.Marshal(m.c)
	if err != nil {
		return
	}
	d2, err := json.Marshal(m.g)
	if err != nil {
		return
	}
	buf.Write(append(d1, '\n'))
	buf.Write(d2)
	m.dumpFile.Truncate(0)
	_, err = m.dumpFile.Write(buf.Bytes())
	if err != nil {
		return
	}
}

func (m *Storage) tickerDump() {
	ticker := time.NewTicker(time.Duration(m.dumpInterval) * time.Second)
	for {
		<-ticker.C
		// fmt.Println(int(t.Second()))
		m.dump()
	}
}

func (m *Storage) readFromFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	content = bytes.ReplaceAll(content, []byte{0}, nil)
	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("invalid file format: expected at least two lines")
	}
	err = json.Unmarshal([]byte(lines[0]), &m.c)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(lines[1]), &m.g)
	if err != nil {
		return err
	}
	return nil
}
