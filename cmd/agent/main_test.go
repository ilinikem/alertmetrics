package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMetric(t *testing.T) {

	type want struct {
		responseCode int
	}

	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{
			name:   "POST 200 OK",
			url:    "http://127.0.0.1:8080/update/gauge/testGauge/100",
			method: http.MethodPost,
			want: want{
				responseCode: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Куда положу ответ
			w := httptest.NewRecorder()
			SendMetric(test.url)
			res := w.Result()

			res.Body.Close()
			if res.StatusCode != test.want.responseCode {
				t.Errorf("got %d, want %d", res.StatusCode, test.want.responseCode)
			}
		})
	}
}
