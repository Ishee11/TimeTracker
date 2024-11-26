package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func GetActivityStatsHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		startDate := c.Query("start_date")             // Получаем дату начала из параметров запроса
		endDate := c.Query("end_date")                 // Получаем дату конца из параметров запроса
		activityTypeFilter := c.Query("activity_type") // Фильтр по типу активности

		var query string       // строковая переменная, в которой будет храниться SQL-запрос
		var args []interface{} // срез, который будет содержать значения для параметров запроса

		/* запрос в базу выглядит следующим образом:
		SELECT type, SUM(duration) as total_duration
		FROM activities
		WHERE user_id = $1
		AND time >= $2
		AND time <= $3
		AND type = $4
		GROUP BY type
		*/

		// Строим SQL-запрос
		query = `SELECT type, SUM(duration) as total_duration
          FROM activities
          WHERE user_id = $1`
		args = append(args, userID)

		/* поскольку фильтры опциональные, мы заранее не знаем сколько будет параметров и используем индекс, который
		меняется, если мы обнаруживаем эти параметры - это позволит избежать ошибки, например, когда мы бы ожидали
		получить дату а вместо нее там был тип активности*/
		paramIndex := 2 // Следующий параметр после userID

		// Добавляем фильтр по датам, если они указаны
		if startDate != "" {
			query += fmt.Sprintf(" AND time >= $%d", paramIndex)
			args = append(args, startDate)
			paramIndex++
		}
		if endDate != "" {
			query += fmt.Sprintf(" AND time <= $%d", paramIndex)
			args = append(args, endDate)
			paramIndex++
		}

		// Добавляем фильтр по типу активности, если он указан
		if activityTypeFilter != "" {
			query += fmt.Sprintf(" AND type = $%d", paramIndex)
			args = append(args, activityTypeFilter)
		}

		// Добавляем GROUP BY для поля type
		query += ` GROUP BY type`

		// отправляем запрос в базу (query - сам запрос, args - параметры $1, $2, $3, $4)
		// получаем объект типа pgx.Rows, который будет содержать данные полученные из базы
		rows, err := db.Query(context.Background(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении статистики"})
			return
		}
		defer rows.Close()

		stats := make(map[string]string) // будет хранить - ключи: типы активности, значения: продолжительность

		for rows.Next() { // циклом проходим по строкам
			var activityType string
			var totalDuration float64
			if err := rows.Scan(&activityType, &totalDuration); err != nil { // извлекаем значения
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных"})
				return
			}

			// Конвертируем продолжительность в удобный формат (часы и минуты)
			hours := int(totalDuration) / 60
			minutes := int(totalDuration) % 60
			if hours > 0 {
				stats[activityType] = fmt.Sprintf("%d часов %d минут", hours, minutes)
			} else {
				stats[activityType] = fmt.Sprintf("%d минут", minutes)
			}
		}

		// Успешный ответ
		c.JSON(http.StatusOK, stats)
	}
}
