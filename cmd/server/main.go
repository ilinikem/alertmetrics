package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ilinikem/alertmetrics/internal/handlers"
	"github.com/ilinikem/alertmetrics/internal/logger"
	"github.com/ilinikem/alertmetrics/internal/middlewares"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func main() {

	// Выполняю парсинг флагов
	parseFlags()

	// Инициализирую хранилище
	memStorage := storage.NewMemStorage()

	// Загружаю данные
	if flagRestore {
		consumer, err := storage.NewConsumer(flagFileStoragePath)
		if err != nil {
			logger.Log.Fatal("Failed to create consumer", zap.Error(err))
		}
		defer consumer.Close()

		loadValue, err := consumer.ReadEvent()
		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				// Файл пустой, инициализируем пустой MemStorage
				logger.Log.Info("File is empty, initializing empty MemStorage.")
				loadValue = &storage.MemStorage{
					Gauge:   make(map[string]storage.Gauge),
					Counter: make(map[string]storage.Counter),
				}
			} else {
				logger.Log.Info("Cannot load values from storage", zap.Error(err))
				// Можно выйти или продолжить с пустыми значениями
				return
			}
		}

		// Загружаю данные только в случае успешного чтения
		memStorage.Counter = loadValue.Counter
		memStorage.Gauge = loadValue.Gauge
		logger.Log.Info("Data loaded from file", zap.String("path", flagFileStoragePath))
	}

	// Создаю Producer для записи данных
	producer, err := storage.NewProducer(flagFileStoragePath)
	if err != nil {
		logger.Log.Fatal("Failed to create producer", zap.Error(err))
	}
	defer producer.Close()

	// Запускаю процесс периодического сохранения метрик в файл
	go func() {
		ticker := time.NewTicker(time.Duration(flagStoreInterval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := producer.WriteEvent(memStorage); err != nil {
				logger.Log.Error("Failed to write metrics to file:", zap.Error(err))
			}
		}
	}()

	// Запускаю сервер
	err = runServer(memStorage)
	if err != nil {
		panic(err)
	}

}

// runServer функция запуска сервера
func runServer(memStorage *storage.MemStorage) error {

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
