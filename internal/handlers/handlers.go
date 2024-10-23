package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ilinikem/alertmetrics/internal/logger"
	"github.com/ilinikem/alertmetrics/internal/models"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

// MetricsHandler Хендлер для работы с метриками
type MetricsHandler struct {
	Storage *storage.MemStorage
}

// NewMetricsHandler Конструктор для создания нового хендлера
func NewMetricsHandler(storage *storage.MemStorage) *MetricsHandler {
	return &MetricsHandler{Storage: storage}
}

func (h *MetricsHandler) UpdateEndpoint(w http.ResponseWriter, r *http.Request) {

	// Разрешаю только POST метод
	if r.Method != http.MethodPost {
		// Логирую запрещенный метод
		logger.Log.Info("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Получаю заголовок Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)

		// Создаю ответ ошибки
		errResp := models.ErrorResponse{
			Message: "Unsupported Media Type",
			Error:   "Header Content-Type must be: 'application/json'",
		}
		resp, err := json.MarshalIndent(errResp, "", "    ")
		if err != nil {
			logger.Log.Info("Error while encoding JSON response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
		return
	}

	var req models.Metrics
	dec := json.NewDecoder(r.Body)

	// Закрываю тело запроса
	defer r.Body.Close()

	if err := dec.Decode(&req); err != nil {
		logger.Log.Info("Error decoding request", zap.Error(err))

		// Создаю ответ ошибки
		errResp := models.ErrorResponse{
			Message: "Error decoding request",
			Error:   err.Error(),
		}
		resp, err := json.MarshalIndent(errResp, "", "    ")
		if err != nil {
			logger.Log.Info("Error while encoding JSON response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return

	}
	switch req.MType {
	case "gauge":
		if req.Value == nil {
			// Создаю ответ ошибки
			errResp := models.ErrorResponse{
				Message: "Value must be provided for gauge",
				Error:   "req.Value = nil",
			}
			resp, err := json.MarshalIndent(errResp, "", "    ")
			if err != nil {
				logger.Log.Info("Error while encoding JSON response", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return

		}
		// Привожу к типу gauge
		g := storage.Gauge(*req.Value)
		h.Storage.UpdateGauge(req.ID, g)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := models.Metrics{
			ID:    req.ID,
			MType: req.MType,
			Value: req.Value,
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			logger.Log.Info("Can not encode response", zap.Error(err))
			return
		}
	case "counter":
		if req.Delta == nil {
			// Создаю ответ ошибки
			errResp := models.ErrorResponse{
				Message: "Delta must be provided for counter",
				Error:   "req.Delta = nil",
			}
			resp, err := json.MarshalIndent(errResp, "", "    ")
			if err != nil {
				logger.Log.Info("Error while encoding JSON response", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
		// Привожу к типу counter
		c := storage.Counter(*req.Delta)
		h.Storage.UpdateCounter(req.ID, c)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := models.Metrics{
			ID:    req.ID,
			MType: req.MType,
			Delta: req.Delta,
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			logger.Log.Info("Can not encode response", zap.Error(err))
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricsHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	// Разрешаю только POST метод
	if r.Method != http.MethodPost {
		// Логирую запрещенный метод
		logger.Log.Info("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Получаю заголовок Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)

		// Создаю ответ ошибки
		errResp := models.ErrorResponse{
			Message: "Unsupported Media Type",
			Error:   "Header Content-Type must be: 'application/json'",
		}
		resp, err := json.MarshalIndent(errResp, "", "    ")
		if err != nil {
			logger.Log.Info("Error while encoding JSON response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
		return
	}

	var req models.Metrics
	dec := json.NewDecoder(r.Body)

	// Закрываю тело запроса
	defer r.Body.Close()

	if err := dec.Decode(&req); err != nil {
		logger.Log.Info("Error decoding request", zap.Error(err))

		// Создаю ответ ошибки
		errResp := models.ErrorResponse{
			Message: "Error decoding request",
			Error:   err.Error(),
		}
		resp, err := json.MarshalIndent(errResp, "", "    ")
		if err != nil {
			logger.Log.Info("Error while encoding JSON response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return

	}

	// Проверяю тип метрики
	switch req.MType {
	case "gauge":
		if value, exists := h.Storage.Gauge[req.ID]; exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			floatValue := float64(value)

			resp := models.Metrics{
				ID:    req.ID,
				MType: req.MType,
				Value: &floatValue,
			}
			enc := json.NewEncoder(w)
			if err := enc.Encode(resp); err != nil {
				logger.Log.Info("Can not encode response", zap.Error(err))
				return
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case "counter":
		if value, exists := h.Storage.Counter[req.ID]; exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			intValue := int64(value)

			resp := models.Metrics{
				ID:    req.ID,
				MType: req.MType,
				Delta: &intValue,
			}
			enc := json.NewEncoder(w)
			if err := enc.Encode(resp); err != nil {
				logger.Log.Info("Can not encode response", zap.Error(err))
				return
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricsHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	allMetrics := make(map[string]string)

	for key, value := range h.Storage.Gauge {
		allMetrics[key] = fmt.Sprintf("%f", value)
	}
	for key, value := range h.Storage.Counter {
		allMetrics[key] = fmt.Sprintf("%d", value)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	htmlResponse := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<title>All metric params</title>
		</head>
		<body>
			%v
		</body>
		</html>
		`, allMetrics)
	_, err := w.Write([]byte(htmlResponse))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
