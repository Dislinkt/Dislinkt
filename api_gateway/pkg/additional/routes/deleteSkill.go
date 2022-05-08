package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteSkill(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	userId := ctx.Param("userId")
	additionId := ctx.Param("additionId")

	res, err := c.DeleteSkill(context.Background(), &pb.EmptyRequest{
		UserId:     userId,
		AdditionId: additionId,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
