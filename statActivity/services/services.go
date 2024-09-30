package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// StatisticsService предоставляет методы для получения статистики активностей
type StatisticsService struct {
	DB *pgxpool.Pool
}

// GetStatistics получает статистику активностей за заданный период
func (s *StatisticsService) GetStatistics(startDate, endDate time.Time, activityType string) (map[string]interface{}, error) {
	ctx := context.Background()

	query := `
		SELECT type, COUNT(*) as count 
		FROM activities
		WHERE time >= $1 AND time <= $2`

	args := []interface{}{startDate, endDate}

	// Добавляем фильтр по типу активности, если он указан
	if activityType != "" {
		query += " AND type = $3"
		args = append(args, activityType)
	}

	query += " GROUP BY type"

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Формируем результат
	result := make(map[string]interface{})
	activityStats := []map[string]interface{}{}
	totalActivities := 0

	for rows.Next() {
		var activityType string
		var count int

		err := rows.Scan(&activityType, &count)
		if err != nil {
			return nil, err
		}

		totalActivities += count
		activityStats = append(activityStats, map[string]interface{}{
			"type":  activityType,
			"count": count,
		})
	}

	result["total_activities"] = totalActivities
	result["activity_stats"] = activityStats

	return result, nil
}

// InitStatisticsService инициализирует StatisticsService
var StatisticsServiceInstance *StatisticsService

func InitStatisticsService(db *pgxpool.Pool) {
	StatisticsServiceInstance = &StatisticsService{DB: db}
}
