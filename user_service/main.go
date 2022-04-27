package main

import (
	"github.com/dislinkt/user-service/startup"
	cfg "github.com/dislinkt/user-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
