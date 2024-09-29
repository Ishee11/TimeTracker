package handlers

import (
	"log"
	"net/http"
	"time-tracker/models"
	"time-tracker/services"

	"github.com/gin-gonic/gin"
)

// AddActivityHandler — это обработчик для добавления активности
func AddActivityHandler(c *gin.Context) {
	var activity models.Activity

	// Декодируем JSON из запроса
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Добавляем логику по добавлению активности в базу данных через ActivityServiceInstance
	if err := services.ActivityServiceInstance.AddActivity(activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении активности"})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusCreated, gin.H{"message": "Активность добавлена"})
	log.Printf("Активность добавлена: %s в %s", activity.Type, activity.Time)
}
