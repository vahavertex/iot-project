package broker

import (
	// "github.com/hooks/auth" // ДОЛЖНО БЫТЬ ТАК
	// "github.com/listeners"  // И ТАК
	"log"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

func SetupBroker() *mqtt.Server {
	s := mqtt.New(&mqtt.Options{
		InlineClient: true, // Добавьте эту строку
	})
	_ = s.AddHook(new(auth.AllowHook), nil)

	tcp := listeners.NewTCP(listeners.Config{
		Type:    listeners.TypeMock,
		Address: "0.0.0.0:1883",
	})
	if err := s.AddListener(tcp); err != nil {
		log.Fatal(err)
	}
	return s
}
