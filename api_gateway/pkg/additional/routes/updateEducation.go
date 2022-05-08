package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateEducation(ctx *gin.Context, c pb.AdditionalUserServiceClient) {
	userId := ctx.Param("userId")
	educationId := ctx.Param("educationId")
	b := NewEducationRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.UpdateEducation(context.Background(), &pb.UpdateEducationRequest{
		UserId:      userId,
		EducationId: educationId,
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

	ctx.JSON(http.StatusOK, &res)
}
