package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

// Для адреса и порта
var flagRunHostAddr string

type Config struct {
	Address string `env:"ADDRESS"`
}

//func parseFlags() {
//
//	flag.StringVar(&flagRunHostAddr, "a", "localhost:8080", "address and port to run server")
//	flag.Parse()
//
//	var cnf Config
//	err := env.Parse(&cnf)
//
//	if err == nil {
//		fmt.Println(cnf.Address)
//		if cnf.Address != "" {
//			flagRunHostAddr = cnf.Address
//		}
//	}
//}

func parseFlags() {

	// Затем парсим флаги
	flag.StringVar(&flagRunHostAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	var cnf Config
	// Сначала парсим переменные окружения
	err := env.Parse(&cnf)
	if err != nil {
		fmt.Println("Ошибка при парсинге переменных окружения:", err)
	}

	// Если переменная окружения была задана, используем её
	if cnf.Address != "" {
		flagRunHostAddr = cnf.Address
	}

	fmt.Println(flagRunHostAddr)

}
