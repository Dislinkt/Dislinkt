package main

import (
	"github.com/dislinkt/message_service/startup"
	cfg "github.com/dislinkt/message_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
