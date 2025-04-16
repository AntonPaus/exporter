package memory

import (
	"bytes"
	"encoding/json"
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
	s := &Storage{
		g:            make(map[string]gauge),
		c:            make(map[string]counter),
		dumpInterval: dumpInterval,
		dumpFile:     nil,
	}
	if restoreFromFile {
		err := s.readFromFile(dumpLocation)
		if err != nil {
			fmt.Println("No storage file found. Continue")
			return nil, err
		}
	}
	dumpFile, err := os.OpenFile(dumpLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	s.dumpFile = dumpFile
	if s.dumpInterval > 0 {
		go s.tickerDump()
	}
	return s, nil
}

func (s *Storage) Terminate() {
	s.dumpFile.Close()
}

func (s *Storage) Update(mName string, mType string, mValue any) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch mType {
	case MetricTypeGauge:
		fmt.Println("----", mValue)
		if val, ok := mValue.(float64); ok {
			s.g[mName] = gauge(val)
		} else {
			return nil, fmt.Errorf("invalid value type for gauge")
		}
	case MetricTypeCounter:
		fmt.Println(mValue)
		if val, ok := mValue.(int64); ok {
			s.c[mName] += counter(val)
		} else {
			return nil, fmt.Errorf("invalid value type for counter")
		}
	default:
		return nil, fmt.Errorf("unknown metric type")
	}
	if s.dumpInterval == 0 {
		s.dump()
	}
	return s.Get(mName, mType)
}

func (s *Storage) Get(mName string, mType string) (any, error) {
	switch mType {
	case MetricTypeGauge:
		if v, ok := s.g[mName]; ok {
			return v, nil
		}
	case MetricTypeCounter:
		if v, ok := s.c[mName]; ok {
			return v, nil
		}
	}
	return nil, fmt.Errorf("metric not found: %s", mName)
}

func (s *Storage) dump() {
	var buf bytes.Buffer
	d1, err := json.Marshal(s.c)
	if err != nil {
		return
	}
	d2, err := json.Marshal(s.g)
	if err != nil {
		return
	}
	buf.Write(append(d1, '\n'))
	buf.Write(d2)
	s.dumpFile.Truncate(0)
	_, err = s.dumpFile.Write(buf.Bytes())
	if err != nil {
		return
	}
}

func (s *Storage) tickerDump() {
	ticker := time.NewTicker(time.Duration(s.dumpInterval) * time.Second)
	for {
		<-ticker.C
		s.dump()
	}
}

func (s *Storage) readFromFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	content = bytes.ReplaceAll(content, []byte{0}, nil)
	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("invalid file format: expected at least two lines")
	}
	err = json.Unmarshal([]byte(lines[0]), &s.c)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(lines[1]), &s.g)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) HealthCheck() error {
	return nil
}
