package main

import (
	"fmt"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"net/http"
	"runtime"
	"time"
)

func main() {

	// Выполняю парсинг флагов
	parseFlags()

	// runAgent Запуск агента с флагами
	runAgent(flagRunHostAddr, flagSendFreq, flagGetFreq)

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
				url := fmt.Sprintf("http://%s/update/gauge/%s/%f", flagRunHostAddr, k, v)
				SendMetric(url)
			}

			for k, v := range memStorage.Counter {
				url := fmt.Sprintf("http://%s/update/counter/%s/%d", flagRunHostAddr, k, v)
				SendMetric(url)
			}

			fmt.Println("Отправляю")
			time.Sleep(time.Duration(flagSendFreq) * time.Second)
		}
	}()
	select {}
}

// SendMetric Отправка метрик
func SendMetric(u string) {
	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

}
