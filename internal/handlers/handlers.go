package handlers

import (
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
