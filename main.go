package main

import (
	"log"
	"net/http"
	"strconv"

	"iot-project/intrernal/broker"
	"iot-project/intrernal/database"
	"iot-project/models"

	"github.com/gorilla/websocket"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

func main() {
	db := database.Init()
	hub := &Hub{clients: make(map[*websocket.Conn]bool), db: db}

	// Инициализация MQTT (используем ваш прошлый код broker.go)
	mqttServer := broker.SetupBroker()

	callbackFn := func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
		val, err := strconv.ParseFloat(string(pk.Payload), 64)
		if err != nil {
			mqttServer.Log.Error("Failed to parse payload", "error", err, "topic", pk.TopicName)
			return
		}

		// Create model object
		reading := models.SensorReading{
			Topic: pk.TopicName,
			Value: val,
		}

		// Save via GORM
		if err := db.Create(&reading).Error; err != nil {
			mqttServer.Log.Error("Failed to save reading", "error", err)
			return
		}

		// Broadcast to browser
		hub.Broadcast(reading)
	}

	// Subscribe to all topics (or specific ones)
	err := mqttServer.Subscribe("#", 1, callbackFn)
	if err != nil {
		log.Fatal(err)
	}

	go mqttServer.Serve()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", hub.handleWS)

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
