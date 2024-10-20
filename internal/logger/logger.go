package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var Log *zap.Logger = zap.NewNop()

type (
	// Структура сведений об ответе
	responseData struct {
		statusCode int
		size       int
	}

	// Добавляю реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Изменяю стандартные интерфейсы
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.statusCode = statusCode
}

func Initialize(level string) error {
	// Парсим уровень логирования
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// Создаю новую конфигурацию логгера
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	// Создаю логер
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

// Создаю middleware-логер для входящих HTTP-запросов
func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			statusCode: 0,
			size:       0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		Log.Info("got request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Duration("duration", duration),
			zap.Int("status_code", responseData.statusCode),
			zap.Int("size", responseData.size),
		)
	})
}
