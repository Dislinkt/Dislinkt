package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/tracer"

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

func (handler *AdditionalUserHandler) GetDegrees(ctx context.Context, request *pb.Get) (*pb.
	GetEntitiesResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetDegreesAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	degrees, err := handler.service.GetDegrees(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.GetEntitiesResponse{
		Entities: mapDegrees(degrees),
	}, nil
}

func (handler *AdditionalUserHandler) GetIndustries(ctx context.Context, request *pb.Get) (*pb.
	GetEntitiesResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetIndustriesAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	industries, err := handler.service.GetIndustries(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.GetEntitiesResponse{
		Entities: mapIndustries(industries),
	}, nil
}

func (handler *AdditionalUserHandler) GetSkills(ctx context.Context, request *pb.Get) (*pb.
	GetEntitiesResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetSkillsAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	skills, err := handler.service.GetSkills(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.GetEntitiesResponse{
		Entities: mapSkills(skills),
	}, nil
}

func (handler *AdditionalUserHandler) GetFieldOfStudies(ctx context.Context, request *pb.Get) (*pb.
	GetEntitiesResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetFieldOfStudies")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fields, err := handler.service.GetFieldOfStudies(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.GetEntitiesResponse{
		Entities: mapFieldsOfStudies(fields),
	}, nil
}

// EDUCATION

func (handler *AdditionalUserHandler) NewEducation(ctx context.Context, request *pb.NewEducationRequest) (*pb.
	EducationResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "NewEducationAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	education := mapNewEducation(request.Education)
	err := handler.service.CreateEducationStart(ctx, request.Id, education)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.EducationResponse{
		Education: mapEducation(education),
	}, nil
}

func (handler *AdditionalUserHandler) GetAllEducation(ctx context.Context, request *pb.GetAllEducationRequest) (*pb.AllEducationResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllEducationAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	educations, err := handler.service.FindUserEducations(ctx, request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllEducationResponse{
		Educations: mapEducations(educations),
	}, nil
}

func (handler *AdditionalUserHandler) UpdateEducation(ctx context.Context, request *pb.UpdateEducationRequest) (*pb.AllEducationResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateEducationAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	educations, err := handler.service.UpdateUserEducationStart(ctx, request.UserId, request.EducationId,
		mapNewEducation(request.Education))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllEducationResponse{
		Educations: mapEducations(educations),
	}, nil
}

func (handler *AdditionalUserHandler) DeleteEducation(ctx context.Context, request *pb.EmptyRequest) (*pb.
	AllEducationResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteEducationAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	educations, err := handler.service.DeleteUserEducationStart(ctx, request.UserId, request.AdditionId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllEducationResponse{
		Educations: mapEducations(educations),
	}, nil
}

// POSITION

func (handler *AdditionalUserHandler) NewPosition(ctx context.Context, request *pb.NewPositionRequest) (*pb.
	PositionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "NewPositionAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	position := mapNewPosition(request.Position)
	createdPosition, err := handler.service.CreatePosition(ctx, request.Id, position)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.PositionResponse{
		Position: mapPosition(createdPosition),
	}, nil
}

func (handler *AdditionalUserHandler) GetAllPosition(ctx context.Context, request *pb.GetAllPositionRequest) (*pb.AllPositionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllPositionAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	positions, err := handler.service.FindUserPositions(ctx, request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllPositionResponse{
		Positions: mapPositions(positions),
	}, nil
}

func (handler *AdditionalUserHandler) UpdatePosition(ctx context.Context, request *pb.UpdatePositionRequest) (*pb.AllPositionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdatePositionAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	positions, err := handler.service.UpdateUserPosition(ctx, request.UserId, request.PositionId,
		mapNewPosition(request.Position))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllPositionResponse{
		Positions: mapPositions(positions),
	}, nil
}

func (handler *AdditionalUserHandler) DeletePosition(ctx context.Context, request *pb.EmptyRequest) (*pb.
	AllPositionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "DeletePositionAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	positions, err := handler.service.DeleteUserPosition(ctx, request.UserId, request.AdditionId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllPositionResponse{
		Positions: mapPositions(positions),
	}, nil
}

// SKILL

func (handler *AdditionalUserHandler) NewSkill(ctx context.Context, request *pb.NewSkillRequest) (*pb.
	SkillResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "NewSkillAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	skill := mapNewSkill(request.Skill)
	err := handler.service.CreateSkillStart(ctx, request.Id, skill)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.SkillResponse{
		Skill: mapSkill(skill),
	}, nil
}

func (handler *AdditionalUserHandler) GetUserSkills(ctx context.Context,
	request *pb.GetUserSkillsRequest) (*pb.UserSkillResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetUserSkillsAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	skills, err := handler.service.FindUserSkills(ctx, request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.UserSkillResponse{
		Skills: mapUserSkills(skills),
	}, nil
}

func (handler *AdditionalUserHandler) UpdateSkill(ctx context.Context, request *pb.UpdateSkillRequest) (*pb.UserSkillResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateSKillAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	skills, err := handler.service.UpdateUserSkillStart(ctx, request.UserId, request.SkillId,
		mapNewSkill(request.Skill))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.UserSkillResponse{
		Skills: mapUserSkills(skills),
	}, nil
}

func (handler *AdditionalUserHandler) DeleteSkill(ctx context.Context, request *pb.EmptyRequest) (*pb.
	UserSkillResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteSkillAPI")
	defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	ctx = tracer.ContextWithSpan(context.Background(), span)
	skills, err := handler.service.DeleteUserSkillStart(ctx, request.UserId, request.AdditionId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.UserSkillResponse{
		Skills: mapUserSkills(skills),
	}, nil
}

// INTEREST

func (handler *AdditionalUserHandler) NewInterest(ctx context.Context, request *pb.NewInterestRequest) (*pb.
	InterestResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "NewInterestAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	interest := mapNewInterest(request.Interest)
	createdInterest, err := handler.service.CreateInterest(ctx, request.Id, interest)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.InterestResponse{
		Interest: mapInterest(createdInterest),
	}, nil
}

func (handler *AdditionalUserHandler) GetAllInterest(ctx context.Context, request *pb.GetAllInterestRequest) (*pb.AllInterestResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllInterestAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	interests, err := handler.service.FindUserInterests(ctx, request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllInterestResponse{
		Interests: mapInterests(interests),
	}, nil
}

func (handler *AdditionalUserHandler) UpdateInterest(ctx context.Context, request *pb.UpdateInterestRequest) (*pb.AllInterestResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateInterestAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	interests, err := handler.service.UpdateUserInterest(ctx, request.UserId, request.InterestId,
		mapNewInterest(request.Interest))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllInterestResponse{
		Interests: mapInterests(interests),
	}, nil
}

func (handler *AdditionalUserHandler) DeleteInterest(ctx context.Context, request *pb.EmptyRequest) (*pb.
	AllInterestResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteInterestAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	interests, err := handler.service.DeleteUserInterest(ctx, request.UserId, request.AdditionId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllInterestResponse{
		Interests: mapInterests(interests),
	}, nil
}
