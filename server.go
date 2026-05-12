package main

import (
	"encoding/json"
	"iot-project/models"
	"net/http"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Hub struct {
	clients map[*websocket.Conn]bool
	db      *gorm.DB
}

func (h *Hub) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	h.clients[conn] = true
	defer func() { delete(h.clients, conn); conn.Close() }()

	// Современный подход: берем последние 50 записей через GORM
	var readings []models.SensorReading
	h.db.Order("created_at desc").Limit(50).Find(&readings)

	// Отправляем историю в правильном порядке (от старых к новым)
	for i := len(readings) - 1; i >= 0; i-- {
		msg, _ := json.Marshal(readings[i])
		conn.WriteMessage(websocket.TextMessage, msg)
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Hub) Broadcast(reading models.SensorReading) {
	msg, _ := json.Marshal(reading)
	for client := range h.clients {
		client.WriteMessage(websocket.TextMessage, msg)
	}
}
