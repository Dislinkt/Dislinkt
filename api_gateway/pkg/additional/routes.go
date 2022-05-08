package additional

import (
	"github.com/dislinkt/api_gateway/pkg/additional/routes"
	"github.com/dislinkt/api_gateway/pkg/auth"
	"github.com/dislinkt/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/additional/user")
	routes.Use(a.AuthRequired)
	routes.POST("/:id/education", svc.NewEducation)
	routes.GET("/:id/education", svc.GetAllEducation)
	routes.PUT("/:userId/education/:educationId", svc.UpdateEducation)
	routes.DELETE("/:userId/education/:additionId", svc.DeleteEducation)

	routes.POST("/:id/position", svc.NewPosition)
	routes.GET("/:id/position", svc.GetAllPosition)
	routes.PUT("/:userId/position/:positionId", svc.UpdatePosition)
	routes.DELETE("/:userId/position/:additionId", svc.DeletePosition)

	routes.POST("/:id/skill", svc.NewSkill)
	routes.GET("/:id/skill", svc.GetAllSkill)
	routes.PUT("/:userId/skill/:skillId", svc.UpdateSkill)
	routes.DELETE("/:userId/skill/:additionId", svc.DeleteSkill)

	routes.POST("/:id/interest", svc.NewInterest)
	routes.GET("/:id/interest", svc.GetAllInterest)
	routes.PUT("/:userId/interest/:interestId", svc.UpdateInterest)
	routes.DELETE("/:userId/interest/:additionId", svc.DeleteInterest)
}

func (svc *ServiceClient) NewEducation(ctx *gin.Context) {
	routes.NewEducation(ctx, svc.Client)
}
func (svc *ServiceClient) GetAllEducation(ctx *gin.Context) {
	routes.GetAllEducation(ctx, svc.Client)
}
func (svc *ServiceClient) UpdateEducation(ctx *gin.Context) {
	routes.UpdateEducation(ctx, svc.Client)
}
func (svc *ServiceClient) DeleteEducation(ctx *gin.Context) {
	routes.DeleteEducation(ctx, svc.Client)
}

func (svc *ServiceClient) NewPosition(ctx *gin.Context) {
	routes.NewPosition(ctx, svc.Client)
}
func (svc *ServiceClient) GetAllPosition(ctx *gin.Context) {
	routes.GetAllPosition(ctx, svc.Client)
}
func (svc *ServiceClient) UpdatePosition(ctx *gin.Context) {
	routes.UpdatePosition(ctx, svc.Client)
}
func (svc *ServiceClient) DeletePosition(ctx *gin.Context) {
	routes.DeletePosition(ctx, svc.Client)
}

func (svc *ServiceClient) NewSkill(ctx *gin.Context) {
	routes.NewSkill(ctx, svc.Client)
}
func (svc *ServiceClient) GetAllSkill(ctx *gin.Context) {
	routes.GetAllSkill(ctx, svc.Client)
}
func (svc *ServiceClient) UpdateSkill(ctx *gin.Context) {
	routes.UpdateSkill(ctx, svc.Client)
}
func (svc *ServiceClient) DeleteSkill(ctx *gin.Context) {
	routes.DeleteSkill(ctx, svc.Client)
}

func (svc *ServiceClient) NewInterest(ctx *gin.Context) {
	routes.NewInterest(ctx, svc.Client)
}
func (svc *ServiceClient) GetAllInterest(ctx *gin.Context) {
	routes.GetAllInterest(ctx, svc.Client)
}
func (svc *ServiceClient) UpdateInterest(ctx *gin.Context) {
	routes.UpdateInterest(ctx, svc.Client)
}
func (svc *ServiceClient) DeleteInterest(ctx *gin.Context) {
	routes.DeleteInterest(ctx, svc.Client)
}
