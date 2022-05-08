package routes

import (
	"context"
	"fmt"
	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterUserRequestBody struct {
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

func RegisterUser(ctx *gin.Context, c pb.UserServiceClient) {
	b := RegisterUserRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	fmt.Println("*********************************************")
	fmt.Println(b.Username)

	res, err := c.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		User: &pb.NewUser{Name: b.Name, Surname: b.Surname, Username: b.Username, Email: b.Email,
			Number: b.Number, Gender: pb.Gender(b.Gender), DateOfBirth: b.DateOfBirth, Password: b.Password,
			Biography: b.Biography, Private: b.Private},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
