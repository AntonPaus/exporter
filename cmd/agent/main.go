package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/AntonPaus/exporter/internal/metrics"
)

func main() {
	address := new(string)
	reportInterval := new(int)
	pollInterval := new(int)
	flag.StringVar(address, "a", "localhost:8080", "server endpoint")
	flag.IntVar(reportInterval, "r", 10, "reportInterval")
	flag.IntVar(pollInterval, "p", 2, "pollInterval")
	flag.Parse()
	osEP := os.Getenv("ADDRESS")
	osRI := os.Getenv("REPORT_INTERVAL")
	osPI := os.Getenv("POLL_INTERVAL")
	if osEP != "" {
		*address = osEP
	}
	if osRI != "" {
		*reportInterval, _ = strconv.Atoi(osRI)
	}
	if osPI != "" {
		*pollInterval, _ = strconv.Atoi(osPI)
	}
	stats := make(chan metrics.Metrics, 60)
	go metrics.PollStats(stats, time.Duration(*pollInterval)*time.Second)
	go metrics.ReportStats(stats, time.Duration(*reportInterval)*time.Second, *address)
	select {}
}
