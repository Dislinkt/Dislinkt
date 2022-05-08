package routes

import (
	"context"
	"net/http"

	pb "github.com/dislinkt/common/proto/auth_service"
	"github.com/gin-gonic/gin"
)

type ValidateRequestBody struct {
	Token string `json:"token"`
}

func Validate(ctx *gin.Context, c pb.AuthServiceClient) {
	b := ValidateRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.ValidateToken(context.Background(), &pb.ValidateRequest{
		Jwt: &pb.JwtToken{Jwt: b.Token},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
