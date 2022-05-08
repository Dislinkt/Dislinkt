package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdatePosition(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	userId := ctx.Param("userId")
	positionId := ctx.Param("positionId")
	b := NewPositionRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.UpdatePosition(context.Background(), &pb.UpdatePositionRequest{
		UserId:     userId,
		PositionId: positionId,
		Position: &pb.NewPosition{
			Title:       b.Title,
			CompanyName: b.CompanyName,
			Industry:    b.Industry,
			StartDate:   b.StartDate,
			EndDate:     b.EndDate,
			Current:     b.Current,
		},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
