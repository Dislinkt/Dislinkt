package main

import (
	"github.com/dislinkt/api_gateway/pkg/auth"
	"github.com/dislinkt/api_gateway/pkg/user"
	cfg "github.com/dislinkt/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config := cfg.NewConfig()
	//server := startup.NewServer(config)
	//server.Start()

	r := gin.Default()

	authSvc := auth.RegisterRoutes(r, config)
	user.RegisterRoutes(r, config, authSvc)

	r.Run(":8000")
}
