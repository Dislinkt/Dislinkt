package routes

import (
	"context"
	"net/http"

	pb "github.com/dislinkt/common/proto/auth_service"
	"github.com/gin-gonic/gin"
)

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context, c pb.AuthServiceClient) {
	b := LoginRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.AuthenticateUser(context.Background(), &pb.LoginRequest{
		UserData: &pb.UserData{Username: b.Username, Password: b.Password},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
