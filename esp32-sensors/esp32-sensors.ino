#include "Config.h"
#include "Network.h"
#include "Sensors.h"

unsigned long lastPublish = 0;
const long interval = 5000; // 5 секунд

void setup() {
  Serial.begin(115200);
  setupWiFi();
  mqttClient.setServer(MQTT_HOST, MQTT_PORT);
}

void loop() {
  if (!mqttClient.connected()) {
    reconnectMQTT();
  }
  mqttClient.loop();

  unsigned long now = millis();
  if (now - lastPublish > interval) {
    lastPublish = now;

    // Модульный обход всех датчиков из конфига
    for (int i = 0; i < sensorCount; i++) {
      float val = readSensorValue(mySensors[i].pin);
      
      // Отправляем данные
      char payload[10];
      dtostrf(val, 4, 2, payload); // Конвертация float в string
      
      mqttClient.publish(mySensors[i].topic, payload);
      
      Serial.printf("Topic: %s | Val: %s\n", mySensors[i].topic, payload);
    }
  }
}
