package handlers

import (
	"fmt"
	"github.com/ilinikem/alertmetrics/internal/storage"
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Получаю ссылку для парсинга
	parsedURL, err := url.Parse(r.URL.String())
	fmt.Println(parsedURL)
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
			_, err := w.Write([]byte(fmt.Sprintf("%.2f", value)))
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
		}
	} else {
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

	var htmlResponse string
	htmlResponse = fmt.Sprintf(`
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
