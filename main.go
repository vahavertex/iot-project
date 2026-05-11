package main

import (
	"log"
	"net/http"
	"strconv"

	"iot-project/models"
	"iot-project/pkg/database"

	"github.com/gorilla/websocket"

	"github.com/mochi-mqtt/server/v2/packets"
)

func main() {
	db := database.Init()
	hub := &Hub{clients: make(map[*websocket.Conn]bool), db: db}

	// Инициализация MQTT (используем ваш прошлый код broker.go)
	mqttServer := setupBroker()

	mqttServer.Events.OnMessage = func(cl *server.Client, pk packets.Packet) {
		val, _ := strconv.ParseFloat(string(pk.Payload), 64)

		// Создаем объект модели
		reading := models.SensorReading{
			Topic: pk.TopicName,
			Value: val,
		}

		// Сохраняем через GORM (одной строкой!)
		db.Create(&reading)

		// Транслируем в браузер
		hub.Broadcast(reading)
	}

	go mqttServer.Serve()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", hub.handleWS)

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
