package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Типы метрик
type gauge float64
type counter int64

// MemStorage Хранение метрик
type MemStorage struct {
	Gauge   map[string]gauge   `json:"gauge"`
	Counter map[string]counter `json:"counter"`
}

// NewMemStorage Конструктор для инициализации MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]gauge),
		Counter: make(map[string]counter),
	}
}

// Storage для хранения метрик
var Storage *MemStorage

// UpdateGauge Метод для обновления Gauge
func (m *MemStorage) UpdateGauge(key string, g gauge) {
	m.Gauge[key] = g
}

// UpdateCounter метод для обновления Counter
func (m *MemStorage) UpdateCounter(key string, c counter) {
	m.Counter[key] += c
}

// Пока не надо
//func mainPage(w http.ResponseWriter, r *http.Request) {
//	// Разрешаю только POST метод
//	if r.Method != http.MethodPost {
//		w.WriteHeader(http.StatusMethodNotAllowed)
//		return
//	}
//	// Устанавливаю заголовок
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//	w.WriteHeader(http.StatusOK)
//
//	_, err := w.Write([]byte("Hello World"))
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}

func updateEndpoint(w http.ResponseWriter, r *http.Request) {
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
		g := gauge(value)
		Storage.UpdateGauge(key, g)

	case "counter":
		// Привожу к int64 и проверяю на ошибки значение
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			return
		}
		// Привожу к типу counter
		c := counter(value)
		Storage.UpdateCounter(key, c)
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

// Вывел для себя полученные метрики
//func getMetrics(w http.ResponseWriter, r *http.Request) {
//
//	// Разрешаю только GET метод
//	if r.Method != http.MethodGet {
//		w.WriteHeader(http.StatusMethodNotAllowed)
//		return
//	}
//
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//	w.WriteHeader(http.StatusOK)
//
//	// Отдаю ответ в формате
//	data, err := json.MarshalIndent(Storage, "", "  ")
//
//	// Если не смог десериализовать данные
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	_, err = w.Write(data)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

func main() {
	// Инициализирую хранилище
	Storage = NewMemStorage()

	mux := http.NewServeMux()
	//mux.HandleFunc("/metrics", getMetrics)
	//mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/update/", updateEndpoint)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
