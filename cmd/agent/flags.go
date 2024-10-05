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
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

func parseFlags() {
	flag.StringVar(&flagRunHostAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagSendFreq, "r", 10, "set frequency for send metrics in seconds")
	flag.IntVar(&flagGetFreq, "p", 2, "set frequency for get metrics in seconds")
	flag.Parse()

	var cnf Config
	err := env.Parse(&cnf)
	if err == nil {
		flagRunHostAddr = cnf.Address
		flagSendFreq = cnf.ReportInterval
		flagGetFreq = cnf.PollInterval
	}

	//// Извлекаю переменную ADDRESS из окружения
	//address := os.Getenv("ADDRESS")
	//
	//if address != "" {
	//	flagRunHostAddr = address
	//}
	//
	//// Извлекаю переменную REPORT_INTERVAL из окружения
	//repInt := os.Getenv("REPORT_INTERVAL")
	//if repInt != "" {
	//	val, err := strconv.Atoi(repInt)
	//	if err != nil {
	//		fmt.Printf("Error parsing %s: %v\n", repInt, err)
	//	}
	//	flagSendFreq = val
	//}
	//
	//// Извлекаю переменную POLL_INTERVAL из окружения
	//polInt := os.Getenv("POLL_INTERVAL")
	//if polInt != "" {
	//	val, err := strconv.Atoi(polInt)
	//	if err != nil {
	//		fmt.Printf("Error parsing %s: %v\n", polInt, err)
	//	}
	//	flagGetFreq = val
	//}
}
