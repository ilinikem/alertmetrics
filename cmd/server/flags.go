package main

import (
	"flag"
	"os"
	"strconv"
)

// Переменные для флагов
var (
	flagRunAddr         string
	flagLogLevel        string
	flagStoreInterval   int
	flagFileStoragePath string
	flagRestore         bool
)

// Парсинг флагов и переменных из окружения
func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "set log level")
	flag.IntVar(&flagStoreInterval, "i", 5, "interval between saving store files")
	flag.StringVar(&flagFileStoragePath, "f", "metrics.json", "path to file storage path")
	flag.BoolVar(&flagRestore, "r", true, "load or not values on start")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envStorageInterval := os.Getenv("STORE_INTERVAL"); envStorageInterval != "" {
		flagStoreInterval, _ = strconv.Atoi(envStorageInterval)
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagFileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if parsedRestore, err := strconv.ParseBool(envRestore); err == nil {
			flagRestore = parsedRestore
		}
	}
}
