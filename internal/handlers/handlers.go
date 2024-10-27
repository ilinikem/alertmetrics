package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ilinikem/alertmetrics/internal/logger"
	"github.com/ilinikem/alertmetrics/internal/models"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	// Получаю ссылку для парсинга
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Разбиваю ссылку на 2 части
	afterUpdate := strings.Split(parsedURL.Path, "/update/")

	// Проверяю, что больше 2-х частей
	if len(afterUpdate) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Разбиваю на сегменты ссылки
	segments := strings.Split(afterUpdate[1], "/")
	if len(segments) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// определяю тип, ключ и значение
	typeKey := segments[0]
	key := segments[1]
	value := segments[2]

	switch typeKey {
	case "gauge":
		// Привожу к float64 и проверяю на ошибки значение
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			return
		}
		// Привожу к типу gauge
		g := storage.Gauge(value)
		h.Storage.UpdateGauge(key, g)

	case "counter":
		// Привожу к int64 и проверяю на ошибки значение
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			return
		}
		// Привожу к типу counter
		c := storage.Counter(value)
		h.Storage.UpdateCounter(key, c)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Устанавливаю заголовок
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("Update page"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *MetricsHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	// Разрешаю только GET метод
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Получаю ссылку для парсинга
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Разбиваю ссылку на 2 части
	afterGet := strings.Split(parsedURL.Path, "/value/")
	// Проверяю, что больше 2-х частей
	if len(afterGet) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Разбиваю на сегменты ссылки
	segments := strings.Split(afterGet[1], "/")

	// определяю тип и имя метрики
	typeMetric := segments[0]
	nameMetric := segments[1]

	if typeMetric == "gauge" {
		if value, exists := h.Storage.Gauge[nameMetric]; exists {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			floatValue := float64(value)
			formattedValue := strconv.FormatFloat(floatValue, 'f', -1, 64)
			_, err := w.Write([]byte(formattedValue))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else if typeMetric == "counter" {
		if value, exists := h.Storage.Counter[nameMetric]; exists {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(fmt.Sprintf("%d", value)))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricsHandler) GetMetricWithJSON(w http.ResponseWriter, r *http.Request) {
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

func (h *MetricsHandler) UpdateEndpointWithJSON(w http.ResponseWriter, r *http.Request) {
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
