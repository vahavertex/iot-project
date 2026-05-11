package models

import (
	"time"
)

type SensorReading struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Topic     string    `gorm:"index" json:"topic"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"` // Автоматически заполнится временем
}
