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

func (handler *UserHandler) GetMe(ctx context.Context, request *pb.GetMeMessage) (*pb.GetMeResponse, error) {
	// span := tracer.StartSpanFromContextMetadata(ctx, "GetAllAPI")
	// defer span.Finish()

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// users, err := handler.service.GetAll(ctx)

	userName := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	fmt.Println("Get me: " + userName)
	user, err := handler.service.FindByUsername(userName)
	if err != nil || user == nil {
		return nil, err
	}
	response := &pb.GetMeResponse{
		User: mapUser(user),
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
	user := mapUpdateUser(request.User)

	// ctx = tracer.ContextWithSpan(context.Background(), span)
	// err := handler.service.Register( ctx, user)
	parsedUUID, err := uuid.FromString(request.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(parsedUUID)
	user.Id = parsedUUID
	dbUser, err := handler.service.StartUpdate(user)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	user.Email = dbUser.Email

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
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	fmt.Println("Patch : " + username)
	user := mapNewUser(request.User)
	user.Username = &username
	err := handler.service.PatchUserStart(user)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.UserResponse{
		User: mapUser(user),
	}, nil
}
