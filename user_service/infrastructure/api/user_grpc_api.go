package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/interceptor"

	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/user_service/application"
	"github.com/dislinkt/user_service/domain"
	uuid "github.com/satori/go.uuid"
)

type UserHandler struct {
	service *application.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.SearchMessage) (*pb.GetAllResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// users, err := handler.service.GetAll(ctx)
	var users *[]domain.User
	var err error
	if len(request.SearchText) == 0 {
		users, err = handler.service.GetAll()
	} else {
		users, err = handler.service.Search(request.SearchText)
	}
	if err != nil || *users == nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range *users {
		current := mapUser(&user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (handler *UserHandler) GetOne(ctx context.Context, request *pb.GetOneMessage) (*pb.UserResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// users, err := handler.service.GetAll(ctx)
	parsedUUID, err := uuid.FromString(request.Id)
	user, err := handler.service.GetOne(parsedUUID)
	if err != nil || user == nil {
		return nil, err
	}
	response := &pb.UserResponse{
		User: mapUser(user),
	}
	return response, nil
}

func (handler *UserHandler) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.
	RegisterUserResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	fmt.Println("Register user")
	user := mapNewUser(request.User)
	fmt.Println("mapper zavrsio")

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	err := handler.service.Register(user)
	fmt.Println("Register zavrsio")
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		User: mapUser(user),
	}, nil
}

func (handler *UserHandler) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.
	UserResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()
	fmt.Println("*************************************************")
	fmt.Println(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	user := mapNewUser(request.User)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	parsedUUID, err := uuid.FromString(request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = handler.service.Update(parsedUUID, user)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.UserResponse{
		User: mapUser(user),
	}, nil
}

func (handler *UserHandler) PatchUser(ctx context.Context, request *pb.PatchUserRequest) (*pb.
	UserResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	parsedUUID, err := uuid.FromString(request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	user, err := handler.service.PatchUser(request.UpdateMask.Paths, mapNewUser(request.User), parsedUUID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.UserResponse{
		User: mapUser(user),
	}, nil
}
