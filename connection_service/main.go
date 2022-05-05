package main

import (
	"github.com/dislinkt/connection_service/startup"
	cfg "github.com/dislinkt/connection_service/startup/config"
)

func main() {

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
