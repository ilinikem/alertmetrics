package main

import (
	"github.com/ilinikem/alertmetrics/internal/handlers"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"net/http"
)

func main() {

	// Инициализирую хранилище
	memStorage := storage.NewMemStorage()
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", metricsHandler.UpdateEndpoint)
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
