Структура проекта

iot-project/
├── models/
│   └── sensor.go    # Модели данных (GORM)
├── pkg/
│   ├── broker/      # Логика MQTT
│   └── database/    # Инициализация БД
├── static/
│   └── index.html
├── main.go
└── server.go
