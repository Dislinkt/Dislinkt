package routes

import (
	"context"
	"net/http"

	pb "github.com/dislinkt/common/proto/post_service"
	"github.com/gin-gonic/gin"
)

func GetAll(ctx *gin.Context, c pb.PostServiceClient) {

	res, err := c.GetAll(context.Background(), &pb.Empty{})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
