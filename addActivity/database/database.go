package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Параметры подключения
var (
	host = "localhost" // "host.docker.internal" //
	//host     = getEnv("DB_HOST", "db")
	port     = getEnvAsInt("DB_PORT", 5432)
	user     = getEnv("DB_USER", "postgres")
	password = getEnv("DB_PASSWORD", "87363699")
	dbname   = getEnv("DB_NAME", "my_time_tracker")
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Вспомогательная функция для чтения числовых переменных окружения
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// ConnectDB устанавливает соединение с базой данных PostgreSQL
func ConnectDB() (*pgxpool.Pool, error) {
	// Формируем строку подключения
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?client_encoding=UTF8", user, password, host, port, dbname)

	// Создаем контекст
	ctx := context.Background()

	// Подключаемся к базе данных PostgreSQL
	connPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных PostgreSQL:", err)
		return nil, err
	}

	// Проверяем соединение с базой данных
	err = connPool.Ping(ctx)
	if err != nil {
		connPool.Close()
		log.Fatal("Не удалось выполнить ping базы данных:", err)
		return nil, err
	}

	fmt.Println("Подключение к базе данных PostgreSQL успешно выполнено")

	return connPool, nil
}
