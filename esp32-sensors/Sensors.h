#include "Config.h"

float readSensorValue(int pin) {
  // Тут может быть dht.readTemperature() или analogRead()
  // Для примера — случайные данные вокруг базового значения
  return (float)analogRead(pin) / 10.0; 
}
