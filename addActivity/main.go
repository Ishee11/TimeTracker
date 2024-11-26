package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time-tracker/database"
	"time-tracker/handlers"
)

func main() {
	router := gin.Default() // создаем сервер

	// Подключение к базе данных
	db, err := database.ConnectDB() // Подключаемся к базе данных
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close() // Закрываем соединение при завершении программы

	// Настраиваем обработчик для добавления активности
	router.POST("/activities", handlers.AddActivityHandler(db))

	// Запуск сервера
	fmt.Println("Сервер запущен на :8080")
	log.Fatal(router.Run("0.0.0.0:8080"))
}
