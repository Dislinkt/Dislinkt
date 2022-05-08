package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewEducationRequestBody struct {
	School       string `json:"school"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

func NewEducation(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	id := ctx.Param("id")
	b := NewEducationRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.NewEducation(context.Background(), &pb.NewEducationRequest{
		Id: id,
		Education: &pb.NewEducation{
			School:       b.School,
			Degree:       b.Degree,
			FieldOfStudy: b.FieldOfStudy,
			StartDate:    b.StartDate,
			EndDate:      b.EndDate,
		},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
