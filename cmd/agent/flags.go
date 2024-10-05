package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

// Для адреса и порта
var flagRunHostAddr string
var flagSendFreq int
var flagGetFreq int

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func parseFlags() {

	flag.StringVar(&flagRunHostAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagSendFreq, "r", 10, "set frequency for send metrics in seconds")
	flag.IntVar(&flagGetFreq, "p", 2, "set frequency for get metrics in seconds")
	flag.Parse()

	var cnf Config
	err := env.Parse(&cnf)
	if err == nil {
		if cnf.Address != "" {
			flagRunHostAddr = cnf.Address
		}
		flagSendFreq = cnf.ReportInterval
		flagGetFreq = cnf.PollInterval
	}
}
