package api

import (
	"context"
	"fmt"

	"github.com/dislinkt/additional_user_service/application"
	pb "github.com/dislinkt/common/proto/additional_user_service"
)

type AdditionalUserHandler struct {
	pb.UnimplementedAdditionalUserServiceServer
	service *application.AdditionalUserService
}

func NewProductHandler(service *application.AdditionalUserService) *AdditionalUserHandler {
	return &AdditionalUserHandler{
		service: service,
	}
}

func (handler *AdditionalUserHandler) NewEducation(ctx context.Context, request *pb.NewEducationRequest) (*pb.
	EducationResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	education := mapNewEducation(request.Education)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Insert( ctx, user)
	createdEducation, err := handler.service.CreateEducation(request.Id, education)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &pb.EducationResponse{
		Education: mapEducation(createdEducation),
	}, nil
}

func (handler *AdditionalUserHandler) GetAllEducation(ctx context.Context, request *pb.GetAllEducationRequest) (*pb.GetAllEducationResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Insert( ctx, user)
	educations, err := handler.service.FindUserEducation(request.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &pb.GetAllEducationResponse{
		Educations: mapEducations(educations),
	}, nil
}
