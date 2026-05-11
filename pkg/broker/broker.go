package main

import (
	"://github.com/hooks/auth" // ДОЛЖНО БЫТЬ ТАК
	"://github.com/listeners"  // И ТАК
	"log"

	"github.com/mochi-mqtt/server/server"
)

func setupBroker() *server.Server {
	s := server.New(nil)
	_ = s.AddHook(new(auth.Allow), nil)

	tcp := listeners.NewTCP("t1", ":1883", nil)
	if err := s.AddListener(tcp); err != nil {
		log.Fatal(err)
	}
	return s
}
