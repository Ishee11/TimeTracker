package main

import (
	"fmt"
	"log"
	"net/http"
	"time-tracker/database"
	"time-tracker/router"
	"time-tracker/services" // Импортируем пакет с обработчиками

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool // Глобальная переменная для хранения подключения к базе данных

func main() {
	// Подключение к базе данных
	var err error
	db, err = database.ConnectDB() // Подключаемся к базе данных
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close() // Закрываем соединение при завершении программы

	// Инициализация сервиса для работы с активностями
	services.InitActivityService(db) // Инициализируем сервис активностей

	// Запуск сервера
	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router.SetupRouter()))
}
