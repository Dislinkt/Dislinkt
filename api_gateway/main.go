package main

import (
	// "github.com/dislinkt/api_gateway/pkg/additional"
	// "github.com/dislinkt/api_gateway/pkg/auth"
	// "github.com/dislinkt/api_gateway/pkg/connection"
	// "github.com/dislinkt/api_gateway/pkg/feed"
	// "github.com/dislinkt/api_gateway/pkg/post"
	// "github.com/dislinkt/api_gateway/pkg/user"
	// "github.com/gin-gonic/gin"

	"github.com/dislinkt/api_gateway/startup"
	cfg "github.com/dislinkt/api_gateway/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
