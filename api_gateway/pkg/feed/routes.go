package feed

import (
	"github.com/dislinkt/api_gateway/infrastructure/api"
	"github.com/dislinkt/api_gateway/pkg/auth"
	"github.com/dislinkt/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)
	userFeedHandler := api.NewUserFeedHandler(c)

	routes := r.Group("/feed")
	routes.Use(a.AuthRequired)
	routes.GET("/user/:id", userFeedHandler.GetUserFeed)
}

//func (userFeedHandler *UserFeedHandler) GetAll(ctx *gin.Context) {
//	routes.GetFeed(ctx, userFeedHandler)
//}
