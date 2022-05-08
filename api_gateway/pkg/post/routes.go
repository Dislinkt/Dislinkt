package post

import (
	"github.com/dislinkt/api_gateway/pkg/auth"
	"github.com/dislinkt/api_gateway/pkg/post/routes"
	"github.com/dislinkt/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/post")
	routes.GET("", svc.getAll)
	routes.GET("/:id", svc.getAllByUserId)
	routes.Use(a.AuthRequired)
	routes.POST("", svc.createPost)
	routes.POST("/:postId/comment", svc.createComment)
	routes.POST("/:postId/like", svc.likePost)
	routes.POST("/:postId/dislike", svc.dislikePost)
}

func (svc *ServiceClient) getAll(ctx *gin.Context) {
	routes.GetAll(ctx, svc.Client)
}
func (svc *ServiceClient) getAllByUserId(ctx *gin.Context) {
	routes.GetAllByUserId(ctx, svc.Client)
}
func (svc *ServiceClient) createPost(ctx *gin.Context) {
	routes.CreatePost(ctx, svc.Client)
}

func (svc *ServiceClient) createComment(ctx *gin.Context) {
	routes.CreateComment(ctx, svc.Client)
}
func (svc *ServiceClient) likePost(ctx *gin.Context) {
	routes.LikePost(ctx, svc.Client)
}
func (svc *ServiceClient) dislikePost(ctx *gin.Context) {
	routes.DislikePost(ctx, svc.Client)
}
