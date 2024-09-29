package main

import (
	"log"
	"net/http"
	"time"
	"time-tracker-statistics/database"
	"time-tracker-statistics/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализируем соединение с базой данных
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	// Инициализируем сервисы
	services.InitStatisticsService(db)

	router := gin.Default()

	// Эндпоинт для получения статистики
	router.GET("/api/statistics", func(c *gin.Context) {
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")
		activityType := c.Query("type")

		startDate, _ := time.Parse("2006-01-02", startDateStr)
		endDate, _ := time.Parse("2006-01-02", endDateStr)

		statistics, err := services.StatisticsServiceInstance.GetStatistics(startDate, endDate, activityType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, statistics)
	})

	router.Run("0.0.0.0:8081")
}
