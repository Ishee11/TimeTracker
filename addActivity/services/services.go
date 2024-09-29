package services

import (
	"context"
	"fmt"
	"log"
	"time"
	"time-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ActivityService предоставляет методы для работы с активностями
type ActivityService struct {
	DB *pgxpool.Pool
}

// AddActivity добавляет новую активность в базу данных
func (s *ActivityService) AddActivity(activity models.Activity) error {
	// Проверяем обязательные поля
	if activity.Type == "" {
		activity.Type = "unknown" // Значение-заглушка для типа активности
	}

	if activity.Time == "" {
		activity.Time = time.Now().Format(time.RFC3339) // Используем текущее время, если не передано
	}

	ctx := context.Background()

	// Попробуем получить последнюю активность
	var lastActivityID uint
	var lastActivityTime string // Используем string для временного поля

	err := s.DB.QueryRow(ctx, `SELECT id, time FROM activities ORDER BY created_at DESC LIMIT 1`).Scan(&lastActivityID, &lastActivityTime)
	if err != nil && err.Error() != "no rows in result set" {
		// Если произошла ошибка (не отсутствие строк), возвращаем её
		log.Printf("Ошибка при запросе последней активности: %v", err)
		return fmt.Errorf("ошибка при запросе последней активности: %w", err)
	}

	// Если есть предыдущая запись, вычисляем продолжительность
	if err == nil {
		// Парсим время последней активности
		parsedLastTime, err := time.Parse(time.RFC3339, lastActivityTime)
		if err != nil {
			log.Printf("Ошибка при парсинге времени последней активности: %v", err)
			return fmt.Errorf("ошибка при парсинге времени последней активности: %w", err)
		}

		currentTime := time.Now()
		duration := currentTime.Sub(parsedLastTime).Minutes() // Время в минутах

		_, err = s.DB.Exec(ctx, `UPDATE activities SET duration = $1 WHERE id = $2`, duration, lastActivityID)
		if err != nil {
			log.Printf("Ошибка при обновлении длительности: %v", err)
			return fmt.Errorf("ошибка при обновлении длительности: %w", err)
		}
	}

	// Выполняем запрос на добавление новой активности
	_, err = s.DB.Exec(ctx, `INSERT INTO activities (type, time, created_at) VALUES ($1, $2, $3)`, activity.Type, activity.Time, time.Now())
	if err != nil {
		log.Printf("Ошибка при добавлении активности: %v", err)
		return fmt.Errorf("ошибка при добавлении активности: %w", err)
	}

	log.Printf("Добавлена активность: %s в %s", activity.Type, activity.Time)
	return nil
}

// InitActivityService инициализирует ActivityService и регистрирует его в глобальной области видимости
var ActivityServiceInstance *ActivityService

func InitActivityService(db *pgxpool.Pool) {
	ActivityServiceInstance = &ActivityService{DB: db}
}
