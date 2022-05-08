package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOne(ctx *gin.Context, c pb.UserServiceClient) {
	id := ctx.Param("id")

	res, err := c.GetOne(context.Background(), &pb.GetOneMessage{
		Id: id,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
