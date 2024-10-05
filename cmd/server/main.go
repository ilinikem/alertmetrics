package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilinikem/alertmetrics/internal/handlers"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"net/http"
)

func main() {

	// Выполняю парсинг флагов
	parseFlags()

	// Запускаю сервер
	err := runServer()
	if err != nil {
		panic(err)
	}

}

// runServer функция запуска сервера
func runServer() error {

	// Инициализирую хранилище
	memStorage := storage.NewMemStorage()
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	// Создаю роутер
	r := chi.NewRouter()
	r.Get("/", metricsHandler.GetAllMetrics)
	r.Get("/value/{typeMetric}/{nameMetric}", metricsHandler.GetMetric)
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", metricsHandler.UpdateEndpoint)

	err := http.ListenAndServe(flagRunHostAddr, r)
	if err != nil {
		return err
	}
	return nil
}
