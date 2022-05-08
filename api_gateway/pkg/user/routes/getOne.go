package routes

import (
	"context"
	"net/http"

	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/gin-gonic/gin"
)

func GetOne(ctx *gin.Context, c pb.UserServiceClient) {

	res, err := c.GetOne(context.Background(), &pb.GetOneMessage{})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
