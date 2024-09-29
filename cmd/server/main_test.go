package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateEndpoint(t *testing.T) {
	// Инициализация NewMemStorage
	Storage = NewMemStorage()

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
			updateEndpoint(w, request)
			res := w.Result()

			// Проверка StatusCode
			assert.Equal(t, test.want.responseCode, res.StatusCode)
			// Проверка заголовка
			assert.Equal(t, test.want.responseHeader, res.Header.Get("content-type"))
		})
	}
}
