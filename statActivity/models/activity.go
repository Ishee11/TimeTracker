package models

import "time"

// Activity представляет собой модель активности
type Activity struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Time      time.Time `json:"time"` // Изменено на time.Time
	CreatedAt time.Time `json:"created_at"`
	Duration  *float64  `json:"duration"` // Указатель для возможности работы с NULL
	UserID    uint      `json:"user_id"`  // Добавлено поле для связи с пользователем
}
