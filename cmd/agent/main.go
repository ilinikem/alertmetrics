package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ilinikem/alertmetrics/internal/models"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"net/http"
	"runtime"
	"time"
)

func main() {

	// Выполняю парсинг флагов
	parseFlags()

	// runAgent Запуск агента с флагами
	runAgent(flagRunAddr, flagSendFreq, flagGetFreq)

}

func runAgent(flagRunHostAddr string, flagSendFreq, flagGetFreq int) {

	memStorage := storage.NewMemStorage()

	go func() {
		for {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			// Сохраняю метрики
			memStorage.UpdateGauge("Alloc", storage.Gauge(memStats.Alloc))
			memStorage.UpdateGauge("BuckHashSys", storage.Gauge(memStats.BuckHashSys))
			memStorage.UpdateGauge("Frees", storage.Gauge(memStats.Frees))
			memStorage.UpdateGauge("GCCPUFraction", storage.Gauge(memStats.GCCPUFraction))
			memStorage.UpdateGauge("GCSys", storage.Gauge(memStats.GCSys))
			memStorage.UpdateGauge("HeapAlloc", storage.Gauge(memStats.HeapAlloc))
			memStorage.UpdateGauge("HeapIdle", storage.Gauge(memStats.HeapIdle))
			memStorage.UpdateGauge("HeapInuse", storage.Gauge(memStats.HeapInuse))
			memStorage.UpdateGauge("HeapObjects", storage.Gauge(memStats.HeapObjects))
			memStorage.UpdateGauge("HeapReleased", storage.Gauge(memStats.HeapReleased))
			memStorage.UpdateGauge("HeapSys", storage.Gauge(memStats.HeapSys))
			memStorage.UpdateGauge("LastGC", storage.Gauge(memStats.LastGC))
			memStorage.UpdateGauge("Lookups", storage.Gauge(memStats.Lookups))
			memStorage.UpdateGauge("MCacheInuse", storage.Gauge(memStats.MCacheInuse))
			memStorage.UpdateGauge("MCacheSys", storage.Gauge(memStats.MCacheSys))
			memStorage.UpdateGauge("MSpanInuse", storage.Gauge(memStats.MSpanInuse))
			memStorage.UpdateGauge("MSpanSys", storage.Gauge(memStats.MSpanSys))
			memStorage.UpdateGauge("Mallocs", storage.Gauge(memStats.Mallocs))
			memStorage.UpdateGauge("NextGC", storage.Gauge(memStats.NextGC))
			memStorage.UpdateGauge("NumForcedGC", storage.Gauge(memStats.NumForcedGC))
			memStorage.UpdateGauge("NumGC", storage.Gauge(memStats.NumGC))
			memStorage.UpdateGauge("OtherSys", storage.Gauge(memStats.OtherSys))
			memStorage.UpdateGauge("PauseTotalNs", storage.Gauge(memStats.PauseTotalNs))
			memStorage.UpdateGauge("StackInuse", storage.Gauge(memStats.StackInuse))
			memStorage.UpdateGauge("StackSys", storage.Gauge(memStats.StackSys))
			memStorage.UpdateGauge("Sys", storage.Gauge(memStats.Sys))
			memStorage.UpdateGauge("TotalAlloc", storage.Gauge(memStats.TotalAlloc))

			memStorage.UpdatePollCount("PollCount")
			memStorage.UpdateRandomValue("RandomValue")

			fmt.Println("Собираю")
			time.Sleep(time.Duration(flagGetFreq) * time.Second)
		}
	}()

	go func() {
		for {
			for k, v := range memStorage.Gauge {
				url := fmt.Sprintf("http://%s/update", flagRunHostAddr)
				val := float64(v)
				metric := models.Metrics{
					ID:    k,
					MType: "gauge",
					Value: &val,
				}
				SendMetric(url, metric)
			}

			for k, v := range memStorage.Counter {
				url := fmt.Sprintf("http://%s/update", flagRunHostAddr)
				delta := int64(v)
				metric := models.Metrics{
					ID:    k,
					MType: "counter",
					Delta: &delta,
				}
				SendMetric(url, metric)
			}

			fmt.Println("Отправляю")
			time.Sleep(time.Duration(flagSendFreq) * time.Second)
		}
	}()
	select {}
}

// SendMetric Отправка метрик
func SendMetric(u string, metric models.Metrics) {
	reqBody, err := json.Marshal(metric)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}

	client := &http.Client{}
	// Добавляю заголовок
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

}
