package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewInterestRequestBody struct {
	Name string `json:"name"`
}

func NewInterest(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	id := ctx.Param("id")
	b := NewInterestRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.NewInterest(context.Background(), &pb.NewInterestRequest{
		Id: id,
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
