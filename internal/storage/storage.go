package storage

import (
	"math/rand"
	"time"
)

type Gauge float64
type Counter int64

// MemStorage Хранение метрик
type MemStorage struct {
	Gauge   map[string]Gauge   `json:"gauge"`
	Counter map[string]Counter `json:"counter"`
}

// NewMemStorage Конструктор для инициализации MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]Gauge),
		Counter: make(map[string]Counter),
	}
}

// UpdateGauge Метод для обновления Gauge
func (m *MemStorage) UpdateGauge(key string, g Gauge) {
	m.Gauge[key] = g
}

// UpdateCounter метод для обновления Counter
func (m *MemStorage) UpdateCounter(key string, c Counter) {
	m.Counter[key] += c
}

// UpdatePollCount увеличивается каждый раз
func (m *MemStorage) UpdatePollCount(key string) {
	m.Counter[key]++
}

// UpdateRandomValue рандомное значение
func (m *MemStorage) UpdateRandomValue(key string) {
	// Инициализирую генератор и сознаю число от 0 до 100
	rand.Seed(time.Now().UnixNano())
	m.Gauge[key] = Gauge(rand.Float64() * 100)
}
