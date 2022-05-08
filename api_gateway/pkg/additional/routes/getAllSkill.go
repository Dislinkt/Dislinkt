package routes

import (
	"context"
	"net/http"

	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
)

func GetAllSkill(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	id := ctx.Param("id")
	res, err := c.GetAllSkill(context.Background(), &pb.GetAllSkillRequest{
		Id: id,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
