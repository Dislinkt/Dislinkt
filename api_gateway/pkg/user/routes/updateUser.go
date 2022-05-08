package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateUser(ctx *gin.Context, c pb.UserServiceClient) {
	id := ctx.Param("id")
	b := RegisterUserRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.UpdateUser(context.Background(), &pb.UpdateUserRequest{
		Id: id,
		User: &pb.NewUser{Name: b.Name, Surname: b.Surname, Username: b.Username, Email: b.Email,
			Number: b.Number, Gender: pb.Gender(b.Gender), DateOfBirth: b.DateOfBirth, Password: b.Password,
			Biography: b.Biography, Private: b.Private},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
