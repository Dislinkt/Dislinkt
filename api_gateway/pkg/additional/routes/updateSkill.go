package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateSkill(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	userId := ctx.Param("userId")
	skillId := ctx.Param("skillId")
	b := NewSkillRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.UpdateSkill(context.Background(), &pb.UpdateSkillRequest{
		UserId:  userId,
		SkillId: skillId,
		Skill: &pb.NewSkill{
			Name: b.Name,
		},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
