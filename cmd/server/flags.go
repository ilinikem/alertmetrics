package main

import "flag"

// Для адреса и порта
var flagRunHostAddr string

func parseFlags() {

	flag.StringVar(&flagRunHostAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}
