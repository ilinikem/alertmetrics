package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilinikem/alertmetrics/internal/handlers"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"net/http"
)

func main() {

	//r := chi.NewRouter()
	//// Получаем машину по id
	//r.Get("/car/{id}", carHandle)
	//
	//// Получаем по brand или по brand и model
	//r.Route("/cars", func(r chi.Router) {
	//	r.Get("/", carsHandle)
	//	r.Route("/{brand}", func(r chi.Router) {
	//		r.Get("/", carsBrandHandle)
	//		r.Get("/{model}", carsBrandModelHandle)
	//	})
	//})
	//log.Fatal(http.ListenAndServe(":8080", r))

	// Инициализирую хранилище
	memStorage := storage.NewMemStorage()
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	//mux := http.NewServeMux()
	//mux.HandleFunc("/update/", metricsHandler.UpdateEndpoint)
	//mux.HandleFunc("/value/", metricsHandler.GetMetric)
	//mux.HandleFunc("/", metricsHandler.GetAllMetrics)
	//err := http.ListenAndServe(":8080", mux)

	r := chi.NewRouter()
	r.Get("/", metricsHandler.GetAllMetrics)
	r.Get("/value/", metricsHandler.GetMetric)
	r.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", metricsHandler.UpdateEndpoint)
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		panic(err)
	}
}
