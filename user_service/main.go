package main

import (
	"github.com/dislinkt/user_service/startup"
	cfg "github.com/dislinkt/user_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
