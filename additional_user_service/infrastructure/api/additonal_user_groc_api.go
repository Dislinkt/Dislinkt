package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/interceptor"
	logger "github.com/dislinkt/common/logging"

	"github.com/dislinkt/additional_user_service/application"
	pb "github.com/dislinkt/common/proto/additional_user_service"
)

type AdditionalUserHandler struct {
	pb.UnimplementedAdditionalUserServiceServer
	service *application.AdditionalUserService
	logger  *logger.Logger
}

func NewProductHandler(service *application.AdditionalUserService) *AdditionalUserHandler {
	logger := logger.InitLogger(context.TODO())
	return &AdditionalUserHandler{
		service: service,
		logger:  logger,
	}
}

func (handler *AdditionalUserHandler) GetDegrees(ctx context.Context, request *pb.Get) (*pb.
	GetEntitiesResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	degrees, err := handler.service.GetDegrees()
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
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	industries, err := handler.service.GetIndustries()
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
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	skills, err := handler.service.GetSkills()
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
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	fields, err := handler.service.GetFieldOfStudies()
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
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("POST rr: EC {%s}", username)
	education := mapNewEducation(request.Education)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	createdEducation, err := handler.service.CreateEducation(request.Id, education)
	if err != nil {
		handler.logger.WarnLogger.Warnf("EC {%s}", username)
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.EducationResponse{
		Education: mapEducation(createdEducation),
	}, nil
}

func (handler *AdditionalUserHandler) GetAllEducation(ctx context.Context, request *pb.GetAllEducationRequest) (*pb.AllEducationResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	educations, err := handler.service.FindUserEducations(request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllEducationResponse{
		Educations: mapEducations(educations),
	}, nil
}

func (handler *AdditionalUserHandler) UpdateEducation(ctx context.Context, request *pb.UpdateEducationRequest) (*pb.AllEducationResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("PUT rr: EU {%s}", username)

	educations, err := handler.service.UpdateUserEducation(request.UserId, request.EducationId,
		mapNewEducation(request.Education))
	if err != nil {
		handler.logger.WarnLogger.Warnf("EU {%s}", username)
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllEducationResponse{
		Educations: mapEducations(educations),
	}, nil
}

func (handler *AdditionalUserHandler) DeleteEducation(ctx context.Context, request *pb.EmptyRequest) (*pb.
	AllEducationResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("DELETE rr: ED {%s}", username)

	educations, err := handler.service.DeleteUserEducation(request.UserId, request.AdditionId)
	if err != nil {
		handler.logger.WarnLogger.Warnf("ED {%s}", username)
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
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("POST rr: PC {%s}", username)
	position := mapNewPosition(request.Position)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	createdPosition, err := handler.service.CreatePosition(request.Id, position)
	if err != nil {
		handler.logger.WarnLogger.Warnf("PC {%s}", username)
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.PositionResponse{
		Position: mapPosition(createdPosition),
	}, nil
}

func (handler *AdditionalUserHandler) GetAllPosition(ctx context.Context, request *pb.GetAllPositionRequest) (*pb.AllPositionResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	positions, err := handler.service.FindUserPositions(request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllPositionResponse{
		Positions: mapPositions(positions),
	}, nil
}

func (handler *AdditionalUserHandler) UpdatePosition(ctx context.Context, request *pb.UpdatePositionRequest) (*pb.AllPositionResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("PUT rr: PU {%s}", username)

	positions, err := handler.service.UpdateUserPosition(request.UserId, request.PositionId,
		mapNewPosition(request.Position))
	if err != nil {
		handler.logger.WarnLogger.Warnf("PU {%s}", username)
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllPositionResponse{
		Positions: mapPositions(positions),
	}, nil
}

func (handler *AdditionalUserHandler) DeletePosition(ctx context.Context, request *pb.EmptyRequest) (*pb.
	AllPositionResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("DELETE rr: PD {%s}", username)

	positions, err := handler.service.DeleteUserPosition(request.UserId, request.AdditionId)
	if err != nil {
		handler.logger.WarnLogger.Warnf("PD {%s}", username)
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
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	skill := mapNewSkill(request.Skill)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	createdSkill, err := handler.service.CreateSkill(request.Id, skill)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.SkillResponse{
		Skill: mapSkill(createdSkill),
	}, nil
}

func (handler *AdditionalUserHandler) GetUserSkills(ctx context.Context,
	request *pb.GetUserSkillsRequest) (*pb.UserSkillResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	skills, err := handler.service.FindUserSkills(request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.UserSkillResponse{
		Skills: mapUserSkills(skills),
	}, nil
}

func (handler *AdditionalUserHandler) UpdateSkill(ctx context.Context, request *pb.UpdateSkillRequest) (*pb.UserSkillResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	skills, err := handler.service.UpdateUserSkill(request.UserId, request.SkillId,
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
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	skills, err := handler.service.DeleteUserSkill(request.UserId, request.AdditionId)
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
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("POST rr: INTC {%s}", username)
	interest := mapNewInterest(request.Interest)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	createdInterest, err := handler.service.CreateInterest(request.Id, interest)
	if err != nil {
		handler.logger.WarnLogger.Warnf("INTC {%s}", username)
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.InterestResponse{
		Interest: mapInterest(createdInterest),
	}, nil
}

func (handler *AdditionalUserHandler) GetAllInterest(ctx context.Context, request *pb.GetAllInterestRequest) (*pb.AllInterestResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	interests, err := handler.service.FindUserInterests(request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllInterestResponse{
		Interests: mapInterests(interests),
	}, nil
}

func (handler *AdditionalUserHandler) UpdateInterest(ctx context.Context, request *pb.UpdateInterestRequest) (*pb.AllInterestResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("PUT rr: INTU {%s}", username)

	interests, err := handler.service.UpdateUserInterest(request.UserId, request.InterestId,
		mapNewInterest(request.Interest))
	if err != nil {
		handler.logger.WarnLogger.Warnf("INTU {%s}", username)
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllInterestResponse{
		Interests: mapInterests(interests),
	}, nil
}

func (handler *AdditionalUserHandler) DeleteInterest(ctx context.Context, request *pb.EmptyRequest) (*pb.
	AllInterestResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("DELETE rr: INTD {%s}", username)

	interests, err := handler.service.DeleteUserInterest(request.UserId, request.AdditionId)
	if err != nil {
		handler.logger.WarnLogger.Warnf("INTD {%s}", username)
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.AllInterestResponse{
		Interests: mapInterests(interests),
	}, nil
}
