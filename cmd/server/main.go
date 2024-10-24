package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilinikem/alertmetrics/internal/handlers"
	"github.com/ilinikem/alertmetrics/internal/logger"
	"github.com/ilinikem/alertmetrics/internal/middlewares"
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

	r.Post("/value", metricsHandler.GetMetricWithJSON)
	r.Post("/update", metricsHandler.UpdateEndpointWithJSON)
	r.Post("/value/", metricsHandler.GetMetricWithJSON)
	r.Post("/update/", metricsHandler.UpdateEndpointWithJSON)
	r.Get("/value/{typeMetric}/{nameMetric}", metricsHandler.GetMetric)
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", metricsHandler.UpdateEndpoint)

	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	err := http.ListenAndServe(flagRunAddr, logger.RequestLogger(middlewares.GzipMiddleware(r)))
	if err != nil {
		logger.Log.Fatal("Error for running server:", zap.Error(err))
	}
	return nil
}
