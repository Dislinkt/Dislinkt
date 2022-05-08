package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateInterest(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	userId := ctx.Param("userId")
	interestId := ctx.Param("interestId")
	b := NewInterestRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.UpdateInterest(context.Background(), &pb.UpdateInterestRequest{
		UserId:     userId,
		InterestId: interestId,
		Interest: &pb.NewInterest{
			Name: b.Name,
		},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
