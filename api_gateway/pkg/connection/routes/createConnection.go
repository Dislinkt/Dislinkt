package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewConnectionRequestBody struct {
	BaseUserUUID    string `json:"base_user_uuid"`
	ConnectUserUUID string `json:"connect_user_uuid"`
}

func CreateConnection(ctx *gin.Context, c pb.ConnectionServiceClient) {
	b := NewConnectionRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.CreateConnection(context.Background(), &pb.NewConnectionRequest{
		Connection: &pb.Connection{BaseUserUUID: b.BaseUserUUID, ConnectUserUUID: b.ConnectUserUUID},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
