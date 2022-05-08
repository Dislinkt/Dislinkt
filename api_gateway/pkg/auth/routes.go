package auth

import (
	"github.com/dislinkt/api_gateway/pkg/auth/routes"
	"github.com/dislinkt/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/api/auth")
	routes.POST("/login", svc.Login)
	routes.POST("/validate", svc.Validate)

	return svc
}

func (svc *ServiceClient) Validate(ctx *gin.Context) {
	routes.Validate(ctx, svc.Client)
}

func (svc *ServiceClient) Login(ctx *gin.Context) {
	routes.Login(ctx, svc.Client)
}
