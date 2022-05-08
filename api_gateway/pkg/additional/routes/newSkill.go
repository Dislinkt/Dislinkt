package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewSkillRequestBody struct {
	Name string `json:"name"`
}

func NewSkill(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	id := ctx.Param("id")
	b := NewSkillRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.NewSkill(context.Background(), &pb.NewSkillRequest{
		Id: id,
		Skill: &pb.NewSkill{
			Name: b.Name,
		},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
