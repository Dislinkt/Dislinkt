package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"net/http"
)

type NewUserRequestBody struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Number      string `json:"number"`
	Gender      int    `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Password    string `json:"password"`
	Biography   string `json:"biography"`
	Private     bool   `json:"private"`
}

type PatchUserRequestBody struct {
	User       NewUserRequestBody `json:"user"`
	UpdateMask []string           `json:"update_mask"`
}

func PatchUser(ctx *gin.Context, c pb.UserServiceClient) {
	id := ctx.Param("id")
	b := PatchUserRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.PatchUser(context.Background(), &pb.PatchUserRequest{
		Id: id,
		User: &pb.NewUser{Name: b.User.Name, Surname: b.User.Surname, Username: b.User.Username, Email: b.User.Email,
			Number: b.User.Number, Gender: pb.Gender(b.User.Gender), DateOfBirth: b.User.DateOfBirth, Password: b.User.Password,
			Biography: b.User.Biography, Private: b.User.Private},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: b.UpdateMask},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
