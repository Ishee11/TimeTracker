package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Параметры подключения
var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "87363699"        // Замените на свой пароль
	dbname   = "my_time_tracker" // Имя базы данных
)

// ConnectDB устанавливает соединение с базой данных PostgreSQL
func ConnectDB() (*pgxpool.Pool, error) {
	// Формируем строку подключения
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)

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

	// Выполняем миграцию модели Activity
	err = migrateActivities(connPool)
	if err != nil {
		log.Fatal("Ошибка миграции модели:", err)
		return nil, err
	}

	return connPool, nil
}

// migrateActivities создает таблицу для модели Activity, если она не существует
func migrateActivities(connPool *pgxpool.Pool) error {
	ctx := context.Background()

	// SQL-запрос для создания таблицы
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS activities (
		id SERIAL PRIMARY KEY,
		type VARCHAR(255) NOT NULL,
		time VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		duration FLOAT
	);`

	// Выполняем запрос
	_, err := connPool.Exec(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу activities: %v", err)
	}

	return nil
}
