package user

import (
	"github.com/dislinkt/api_gateway/pkg/auth"
	"github.com/dislinkt/api_gateway/pkg/user/routes"
	"github.com/dislinkt/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/user")
	routes.POST("", svc.RegisterUser)
	routes.GET("", svc.GetAll)
	routes.GET("/:id", svc.GetOne)
	routes.Use(a.AuthRequired)
	{
		routes.PUT("/:id", a.Authorize, svc.UpdateUser)
		routes.PATCH("/:id", a.Authorize, svc.PatchUser)
	}

}

func (svc *ServiceClient) GetAll(ctx *gin.Context) {
	routes.GetAll(ctx, svc.Client)
}
func (svc *ServiceClient) GetOne(ctx *gin.Context) {
	routes.GetOne(ctx, svc.Client)
}
func (svc *ServiceClient) RegisterUser(ctx *gin.Context) {
	routes.RegisterUser(ctx, svc.Client)
}

func (svc *ServiceClient) UpdateUser(ctx *gin.Context) {
	routes.UpdateUser(ctx, svc.Client)
}
func (svc *ServiceClient) PatchUser(ctx *gin.Context) {
	routes.PatchUser(ctx, svc.Client)
}
