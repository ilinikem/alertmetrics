package main

import "flag"

// Для адреса и порта
var flagRunHostAddr string
var flagSendFreq int
var flagGetFreq int

func parseFlags() {
	flag.StringVar(&flagRunHostAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagSendFreq, "r", 10, "set frequency for send metrics in seconds")
	flag.IntVar(&flagGetFreq, "p", 2, "set frequency for get metrics in seconds")
	flag.Parse()
}
