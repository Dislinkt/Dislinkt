package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterUserRequestBody struct {
	user NewUserRequestBody `json:"user"`
}

func RegisterUser(ctx *gin.Context, c pb.UserServiceClient) {
	b := RegisterUserRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		User: &pb.NewUser{Name: b.user.name, Surname: b.user.surname, Username: b.user.username, Email: b.user.email,
			Number: b.user.number, Gender: pb.Gender(b.user.gender), DateOfBirth: b.user.date_of_birth, Password: b.user.password,
			Biography: b.user.biography, Private: b.user.private},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
