#ifndef CONFIG_H
#define CONFIG_H

// Настройки Wi-Fi
#define WIFI_SSID "vaha"
#define WIFI_PASS "Qw123456"

// Настройки Go-сервера
#define MQTT_HOST IPAddress(192, 168, 1, 50) 
#define MQTT_PORT 1883

// Описание датчика
struct Sensor {
  const char* topic;
  int pin;
  float lastValue;
};

// Массив ваших датчиков (N штук)
Sensor mySensors[] = {
  {"home/livingroom/temp", 32, 0.0},
  {"home/livingroom/hum",  33, 0.0},
  {"home/kitchen/gas",     34, 0.0}
};

const int sensorCount = sizeof(mySensors) / sizeof(Sensor);

#endif
