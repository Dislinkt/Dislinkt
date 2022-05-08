package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/post_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateComment(ctx *gin.Context, c pb.PostServiceClient) {
	postId := ctx.Param("postId")
	b := Comment{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.CreateComment(context.Background(), &pb.CreateCommentRequest{
		PostId:  postId,
		Comment: &pb.Comment{Username: b.Username, CommentText: b.CommentText},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
