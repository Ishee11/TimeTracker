package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"users/database"
	"users/handlers"
)

func main() {
	// Инициализируем соединение с базой данных
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	router := gin.Default()

	router.POST("/users/create", handlers.CreateUserHandler(db))
	router.POST("/users/login", handlers.LoginHandler(db))

	// Запускаем сервер
	router.Run("0.0.0.0:8082")
}
