package main

import (
	"github.com/dislinkt/additional_user_service/startup"
	cfg "github.com/dislinkt/additional_user_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
