package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/ilinikem/alertmetrics/internal/handlers"
	"github.com/ilinikem/alertmetrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateEndpoint(t *testing.T) {
	// Инициализация NewMemStorage
	memStorage := storage.NewMemStorage()
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	type want struct {
		responseCode   int
		responseHeader string
	}
	tests := []struct {
		name     string
		endpoint string
		method   string
		want     want
	}{
		{
			name:     "Not POST method",
			method:   http.MethodGet,
			endpoint: "/update",
			want: want{
				responseCode:   http.StatusMethodNotAllowed,
				responseHeader: "",
			},
		},
		{
			name:     "POST method",
			method:   http.MethodPost,
			endpoint: "/update/gauge/testGauges/22.5",
			want: want{
				responseCode:   http.StatusOK,
				responseHeader: "text/plain; charset=utf-8",
			},
		},
		{
			name:     "POST method with unknown metric type",
			method:   http.MethodPost,
			endpoint: "/update/gaugeTest/testGauges/22.5",
			want: want{
				responseCode:   http.StatusBadRequest,
				responseHeader: "",
			},
		},
		{
			name:     "POST method with bad metric value",
			method:   http.MethodPost,
			endpoint: "/update/gaugeTest/testGauges/testValue",
			want: want{
				responseCode:   http.StatusBadRequest,
				responseHeader: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//Создаю запрос
			request := httptest.NewRequest(test.method, test.endpoint, nil)
			// Куда положу ответ
			w := httptest.NewRecorder()
			metricsHandler.UpdateEndpoint(w, request)
			res := w.Result()

			// Закрываю тело ответа
			res.Body.Close()

			// Проверка StatusCode
			assert.Equal(t, test.want.responseCode, res.StatusCode)
			// Проверка заголовка
			assert.Equal(t, test.want.responseHeader, res.Header.Get("content-type"))
		})
	}
}

func TestUpdateEndpointWithJSON(t *testing.T) {
	// Создаю объект хранилища
	memStorage := storage.NewMemStorage()
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	// Создаю роутер
	r := chi.NewRouter()
	r.Post("/update", metricsHandler.UpdateEndpointWithJSON)

	// Создаю тестовый сервер
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		name         string
		url          string
		method       string
		header       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "method_get",
			url:          "/update",
			method:       http.MethodGet,
			header:       "application/json",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_delete",
			url:          "/update",
			method:       http.MethodDelete,
			header:       "application/json",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_put",
			url:          "/update",
			method:       http.MethodPut,
			header:       "application/json",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "post_request_with_wrong_header",
			url:          "/update",
			method:       http.MethodPost,
			header:       "text/plain",
			body:         "",
			expectedCode: http.StatusUnsupportedMediaType,
			expectedBody: "",
		},
		{
			name:         "method_post_update_counter_success",
			url:          "/update",
			method:       http.MethodPost,
			header:       "application/json",
			body:         `{"id": "PollCount", "type": "counter", "delta": 12}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id": "PollCount", "type": "counter", "delta": 12}`,
		},
		{
			name:         "method_post_update_gauge_success",
			url:          "/update",
			method:       http.MethodPost,
			header:       "application/json",
			body:         `{"id": "Alloc", "type": "gauge", "value": 34.1}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id": "Alloc", "type": "gauge", "value": 34.1}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.SetHeader("Content-Type", tc.header)
			req.URL = srv.URL + tc.url

			if len(tc.body) > 0 {
				req.SetBody(tc.body)
			}
			resp, err := req.Send()
			assert.NoError(t, err, "got error while making request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "got status code")
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}
