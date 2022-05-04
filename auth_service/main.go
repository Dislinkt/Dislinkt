package main

import (
	"github.com/dislinkt/auth-service/startup"
	cfg "github.com/dislinkt/auth-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
