package api

import (
	"context"
	"fmt"

	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/user-service/application"
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

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.EmptyMessage) (*pb.GetAllResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// users, err := handler.service.GetAll(ctx)
	users, err := handler.service.GetAll()
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

func (handler *UserHandler) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.
	RegisterUserResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	user := mapNewUser(request.User)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Insert( ctx, user)
	err := handler.service.Insert(user)
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

	user := mapNewUser(request.User)
	fmt.Println(user.Biography)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Insert( ctx, user)
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
	// err := handler.service.Insert( ctx, user)
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
