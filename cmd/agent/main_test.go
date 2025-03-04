package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_pollStats(t *testing.T) {
	stats := make(chan Metrics, 1)
	go pollStats(stats, 2*time.Second)
	time.Sleep(4 * time.Second)
	receivedStats := <-stats
	assert.NotZero(t, receivedStats.Alloc)
	assert.NotZero(t, receivedStats.Frees)
	assert.NotZero(t, receivedStats.TotalAlloc)
	assert.NotZero(t, receivedStats.PollCount)
}
