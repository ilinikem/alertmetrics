package main

import (
	"flag"
	"os"
	"strconv"
)

// Для флагов
var (
	flagRunAddr  string
	flagSendFreq int
	flagGetFreq  int
	flagLogLevel string
)

func parseFlags() {

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagSendFreq, "r", 10, "set frequency for send metrics in seconds")
	flag.IntVar(&flagGetFreq, "p", 2, "set frequency for get metrics in seconds")
	flag.StringVar(&flagLogLevel, "l", "info", "set log level")

	flag.Parse()

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envSendFreq := os.Getenv("REPORT_INTERVAL"); envSendFreq != "" {
		if value, err := strconv.Atoi(envSendFreq); err == nil {
			flagSendFreq = value
		}
	}
	if envGetFreq := os.Getenv("POLL_INTERVAL"); envGetFreq != "" {
		if value, err := strconv.Atoi(envGetFreq); err == nil {
			flagGetFreq = value
		}
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}

}
