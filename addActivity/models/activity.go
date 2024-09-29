package models

// Activity представляет собой модель активности
type Activity struct {
	ID        uint    `json:"id" gorm:"primaryKey"` // Уникальный идентификатор активности
	Type      string  `json:"type"`                 // Тип активности
	Time      string  `json:"time"`                 // Дата и время активности в формате ISO 8601
	CreatedAt string  `json:"created_at"`           // Дата и время создания записи
	Duration  float64 `json:"duration"`             // Продолжительность активности в часах
}
