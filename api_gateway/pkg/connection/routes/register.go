package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequestBody struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

func Register(ctx *gin.Context, c pb.ConnectionServiceClient) {
	b := RegisterRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.Register(context.Background(), &pb.RegisterRequest{
		User: &pb.User{UserID: b.UserID, Status: b.Status},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
