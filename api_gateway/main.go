package main

import (
	"github.com/dislinkt/api_gateway/startup"
	cfg "github.com/dislinkt/api_gateway/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
