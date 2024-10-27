package storage

import (
	"encoding/json"
	"math/rand"
	"os"
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
	rand.New(rand.NewSource(time.Now().UnixNano()))
	m.Gauge[key] = Gauge(rand.Float64() * 100)
}

// Producer записывает в файл
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer создает Producer
func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteEvent записывает состояние MemStorage в файл
func (p *Producer) WriteEvent(memStorage *MemStorage) error {
	// Закрываем текущий файл
	if err := p.file.Close(); err != nil {
		return err
	}

	// Открываем его снова с флагом os.O_TRUNC для перезаписи
	file, err := os.OpenFile(p.file.Name(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	// Обновляем файл и encoder
	p.file = file
	p.encoder = json.NewEncoder(file)
	// Записываем новые данные
	return p.encoder.Encode(memStorage)
}

// Close закрывает файл Producer
func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*MemStorage, error) {
	event := &MemStorage{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
