package storage

import (
	"context"
	"fmt"
	"github.com/ilinikem/alertmetrics/internal/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// Переменная для подключения к базе данных
var db *pgx.Conn

// Функция инициализации подключения к базе данных
func InitDB(dsn string) error {
	var err error
	// Подключаемся к базе данных
	db, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		// Логируем ошибку подключения
		logger.Log.Error("Error while connecting to the database", zap.Error(err))
		return fmt.Errorf("error while connecting to the database: %v", err)
	}

	// Проверка соединения с базой данных
	err = db.Ping(context.Background())
	if err != nil {
		// Логируем ошибку пинга
		logger.Log.Error("Failed to establish connection with the database", zap.Error(err))
		return fmt.Errorf("failed to establish connection with the database: %v", err)
	}

	return nil
}

// Функция проверки соединения с базой данных
func PingDB() error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return db.Ping(context.Background())
}
