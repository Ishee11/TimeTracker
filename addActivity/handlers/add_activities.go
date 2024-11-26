package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
	"time-tracker/models"
)

// создаем обработчик для добавления активности
func AddActivityHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var activity models.Activity

		// Декодируем JSON из запроса
		err := c.ShouldBindJSON(&activity) // преобразуем полученные данные в структуру models.Activity
		if err != nil {                    // если не получилось:
			log.Printf("Ошибка при декодировании JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
			return
		}

		// Логируем входящие данные для отладки
		log.Printf("Попытка добавить активность: %+v", activity)

		ctx := context.Background() // создаем пустой контекст (его требует db.Exec как параметр)

		// Выполняем запрос на добавление новой активности
		_, err = db.Exec(ctx, `INSERT INTO activities (user_id, type, time, created_at) VALUES ($1, $2, $3, $4)`,
			activity.UserID, activity.Type, activity.Time, time.Now())
		if err != nil {
			log.Printf("Ошибка при добавлении активности: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении активности", "details": err.Error()})
			return
		}

		// Успешный ответ
		c.JSON(http.StatusCreated, gin.H{"message": "Активность добавлена"})
		log.Printf("Активность добавлена: %s в %s, userID: %v", activity.Type, activity.Time, activity.UserID)
	}
}
