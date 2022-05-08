package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewPositionRequestBody struct {
	Title       string `json:"title"`
	CompanyName string `json:"company_name"`
	Industry    string `json:"industry"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Current     bool   `json:"current"`
}

func NewPosition(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	id := ctx.Param("id")
	b := NewPositionRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.NewPosition(context.Background(), &pb.NewPositionRequest{
		Id: id,
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

	ctx.JSON(http.StatusCreated, &res)
}
