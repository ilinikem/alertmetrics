package main

import (
	"github.com/ilinikem/alertmetrics/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMetric(t *testing.T) {
	value := 33.7

	type want struct {
		responseCode int
	}

	tests := []struct {
		name   string
		url    string
		method string
		metric models.Metrics
		want   want
	}{
		{
			name:   "POST 200 OK",
			url:    "http://127.0.0.1:8080/update",
			method: http.MethodPost,
			metric: models.Metrics{
				ID:    "Alloc",
				MType: "gauge",
				Value: &value,
			},
			want: want{
				responseCode: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Куда положу ответ
			w := httptest.NewRecorder()
			SendMetric(test.url, test.metric)
			res := w.Result()

			res.Body.Close()
			if res.StatusCode != test.want.responseCode {
				t.Errorf("got %d, want %d", res.StatusCode, test.want.responseCode)
			}
		})
	}
}
