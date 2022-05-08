package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"net/http"
)

type NewUserRequestBody struct {
	name          string `json:"name"`
	surname       string `json:"surname"`
	username      string `json:"username"`
	email         string `json:"email"`
	number        string `json:"number"`
	gender        int    `json:"gender"`
	date_of_birth string `json:"date_of_birth"`
	password      string `json:"password"`
	biography     string `json:"biography"`
	private       bool   `json:"private"`
}

type PatchUserRequestBody struct {
	id          string             `json:"id"`
	user        NewUserRequestBody `json:"user"`
	update_mask []string           `json:"update_mask"`
}

func PatchUser(ctx *gin.Context, c pb.UserServiceClient) {
	b := PatchUserRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.PatchUser(context.Background(), &pb.PatchUserRequest{
		Id: b.id,
		User: &pb.NewUser{Name: b.user.name, Surname: b.user.surname, Username: b.user.username, Email: b.user.email,
			Number: b.user.number, Gender: pb.Gender(b.user.gender), DateOfBirth: b.user.date_of_birth, Password: b.user.password,
			Biography: b.user.biography, Private: b.user.private},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: b.update_mask},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
