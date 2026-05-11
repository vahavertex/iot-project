#include <WiFi.h>
#include <PubSubClient.h>
#include "Config.h"

WiFiClient espClient;
PubSubClient mqttClient(espClient);

void setupWiFi() {
  WiFi.begin(WIFI_SSID, WIFI_PASS);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500); Serial.print(".");
  }
  Serial.println("\nWiFi Connected");
}

void reconnectMQTT() {
  while (!mqttClient.connected()) {
    if (mqttClient.connect("ESP32_Client")) {
      Serial.println("MQTT Connected");
    } else {
      delay(5000);
    }
  }
}
