package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time-tracker-statistics/database"
	"time-tracker-statistics/handlers"
)

func main() {
	// Инициализируем соединение с базой данных
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/activity/stats", handlers.GetActivityStatsHandler(db))

	// Запускаем сервер
	router.Run(":8081")
}
