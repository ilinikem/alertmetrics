package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

// Для адреса и порта
var flagRunHostAddr string

type Config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func parseFlags() {

	flag.StringVar(&flagRunHostAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	var cnf Config
	err := env.Parse(&cnf)
	if err == nil {
		if cnf.Address != "" {
			flagRunHostAddr = cnf.Address
		}
	}
}
