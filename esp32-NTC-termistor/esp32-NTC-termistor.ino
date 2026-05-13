#include <WiFi.h>
#include <PubSubClient.h>
#include <cmath>

// Настройки Wi-Fi
const char* ssid = "CIR";
const char* password = "CIR2024GGNTU";

// Настройки MQTT-брокера
const char* mqtt_server = "192.168.16.91";  // Например, "111.222.3.4" или "hivemq.com"
const int mqtt_port = 1883;
const char* mqtt_topic = "esp32/sensor/temperature";

// Настройки термистора
const int thermistorPin = 34;           // GPIO 34 (ADC1)
const float SERIES_RESISTOR = 10000.0;  // Резистор делителя 10 кОм
const float B_COEFFICIENT = 3950.0;     // B-параметр термистора
const float THERMISTOR_NOMINAL = 10000.0;
const float TEMPERATURE_NOMINAL = 25.0;

WiFiClient espClient;
PubSubClient client(espClient);
unsigned long lastMsg = 0;

void setup_wifi() {
  delay(10);
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
  }
}

void reconnect() {
  while (!client.connected()) {
    String clientId = "ESP32Client-";
    clientId += String(random(0, 0xffff), HEX);
    if (client.connect(clientId.c_str())) {
      // Подключение успешно
    } else {
      delay(5000);
    }
  }
}

void setup() {
  Serial.begin(115200);
  setup_wifi();
  client.setServer(mqtt_server, mqtt_port);

  // Настройка АЦП: 12 бит разрешения (0-4095)
  analogReadResolution(12);
}

void loop() {
  if (!client.connected()) {
    reconnect();
  }
  client.loop();

  unsigned long now = millis();
  if (now - lastMsg > 5000) {  // Отправка каждые 5 секунд
    lastMsg = now;

    float analogValue = analogRead(thermistorPin);    

    // Исключаем деление на ноль при максимальном значении АЦП
    if (analogValue >= 4095) analogValue = 4094;
    if (analogValue <= 0) analogValue = 1;

    // Расчет сопротивления для АЦП 12-бит (4095)
    float resistance = SERIES_RESISTOR / (4095.0 / analogValue - 1.0);

    /*
    // Вместо: float analogValue = analogRead(thermistorPin);
    // Используем получение калиброванных милливольт:
      float v_out = analogReadMilliVolts(thermistorPin); 

    // Исключаем системные ошибки (чтобы не делить на 0)
      if (v_out >= 3300.0) v_out = 3299.0;
      if (v_out <= 0.0) v_out = 1.0;

    // Расчет сопротивления термистора на основе милливольт (схема делителя к GND):
    // R = R_постоянный / ((V_питания / V_выхода) - 1)
      float resistance = SERIES_RESISTOR / ((3300.0 / v_out) - 1.0);
    */

    // Формула Стейнхарта — Харта
    float steinhart;
    steinhart = resistance / THERMISTOR_NOMINAL;
    steinhart = log(steinhart);
    steinhart /= B_COEFFICIENT;
    steinhart += 1.0 / (TEMPERATURE_NOMINAL + 273.15);
    steinhart = 1.0 / steinhart;
    steinhart -= 273.15;

    // Публикация в MQTT
    String tempStr = String(steinhart, 2);
    client.publish(mqtt_topic, tempStr.c_str());

    Serial.print("Отправлено на MQTT: ");
    Serial.println(tempStr);
  }
}
