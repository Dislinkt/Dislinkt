package routes

import (
	"context"
	"net/http"

	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/gin-gonic/gin"
)

func GetAllConnectionForUser(ctx *gin.Context, c pb.ConnectionServiceClient) {
	id := ctx.Param("id")
	res, err := c.GetAllConnectionForUser(context.Background(), &pb.GetConnectionRequest{Uuid: id})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
