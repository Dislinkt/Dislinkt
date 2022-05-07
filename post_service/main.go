package main

import (
	"post_service/startup"
	cfg "post_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
