package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/post_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReactionRequest struct {
	Username string `json:"username"`
}

func LikePost(ctx *gin.Context, c pb.PostServiceClient) {
	postId := ctx.Param("postId")
	b := ReactionRequest{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.LikePost(context.Background(), &pb.ReactionRequest{
		PostId:   postId,
		Username: b.Username,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
