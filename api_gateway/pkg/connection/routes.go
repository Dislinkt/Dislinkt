package connection

import (
	"github.com/dislinkt/api_gateway/pkg/auth"
	"github.com/dislinkt/api_gateway/pkg/connection/routes"
	"github.com/dislinkt/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/connection")
	routes.POST("/register", svc.Register)
	routes.GET("/user/:userId", svc.GetAllConnectionForUser)
	routes.Use(a.AuthRequired)
	routes.POST("", svc.CreateConnection)
	routes.PUT("/accept", svc.AcceptConnection)
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	routes.Register(ctx, svc.Client)
}
func (svc *ServiceClient) GetAllConnectionForUser(ctx *gin.Context) {
	routes.GetAllConnectionForUser(ctx, svc.Client)
}
func (svc *ServiceClient) CreateConnection(ctx *gin.Context) {
	routes.CreateConnection(ctx, svc.Client)
}

func (svc *ServiceClient) AcceptConnection(ctx *gin.Context) {
	routes.AcceptConnection(ctx, svc.Client)
}
