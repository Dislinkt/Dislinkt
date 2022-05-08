package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/post_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllByUserId(ctx *gin.Context, c pb.PostServiceClient) {
	id := ctx.Param("id")

	res, err := c.GetAllByUserId(context.Background(), &pb.GetRequest{
		Id: id,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
