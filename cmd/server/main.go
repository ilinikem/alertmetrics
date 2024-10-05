package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ilinikem/alertmetrics/internal/handlers"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"net/http"
)

func main() {

	// Выполняю парсинг флагов
	parseFlags()

	// Запускаю сервер
	err := runServer(flagRunHostAddr)
	if err != nil {
		panic(err)
	}

}

// runServer функция запуска сервера
func runServer(flagRunHostAddr string) error {

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
		fmt.Println("Ошибка запуска сервера:", err)
		return err
	}
	return nil
}
