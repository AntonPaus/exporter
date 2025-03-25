package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPollStats(t *testing.T) {
	m := NewMetrics()
	go m.Poll(time.Duration(1 * time.Second))
	time.Sleep(2 * time.Second)
	fmt.Println(m.values)
	assert.NotZero(t, m.values.Alloc)
	assert.NotZero(t, m.values.Frees)
	assert.NotZero(t, m.values.TotalAlloc)
	assert.NotZero(t, m.values.PollCount)
}
