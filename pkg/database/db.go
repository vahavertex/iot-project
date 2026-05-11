package database

import (
	"iot-project/models"
	"log"

	"gorm.io/gorm"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("sensors.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// Авто-миграция (создаст/обновит таблицы согласно структурам)
	db.AutoMigrate(&models.SensorReading{})
	return db
}
