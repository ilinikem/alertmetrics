package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ilinikem/alertmetrics/internal/handlers"
	"github.com/ilinikem/alertmetrics/internal/logger"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"go.uber.org/zap"
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

	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	err := http.ListenAndServe(flagRunAddr, logger.RequestLogger(r))

	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return err
	}
	return nil
}
