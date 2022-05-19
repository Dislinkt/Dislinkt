package main

import (
	//"github.com/dislinkt/api_gateway/pkg/additional"
	//"github.com/dislinkt/api_gateway/pkg/auth"
	//"github.com/dislinkt/api_gateway/pkg/connection"
	//"github.com/dislinkt/api_gateway/pkg/feed"
	//"github.com/dislinkt/api_gateway/pkg/post"
	//"github.com/dislinkt/api_gateway/pkg/user"
	//"github.com/gin-gonic/gin"

	"github.com/dislinkt/api_gateway/startup"
	cfg "github.com/dislinkt/api_gateway/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()

	//r := gin.Default()
	//
	//authSvc := auth.RegisterRoutes(r, config)
	//user.RegisterRoutes(r, config, authSvc)
	//additional.RegisterRoutes(r, config, authSvc)
	//connection.RegisterRoutes(r, config, authSvc)
	//post.RegisterRoutes(r, config, authSvc)
	//feed.RegisterRoutes(r, config, authSvc)
	//
	//r.Run(":8000")
}
