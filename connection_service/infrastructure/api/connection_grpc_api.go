package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/interceptor"
	logger "github.com/dislinkt/common/logging"

	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/application"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	service *application.ConnectionService
	logger  *logger.Logger
}

func NewConnectionHandler(service *application.ConnectionService) *ConnectionHandler {
	logger := logger.InitLogger(context.TODO())
	return &ConnectionHandler{
		service: service,
		logger:  logger,
	}
}

func (handler *ConnectionHandler) Register(ctx context.Context, request *pb.RegisterRequest) (response *pb.RegisterResponse, err error) {
	fmt.Println("[ConnectionHandler]:Register")

	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("POST rr: UC {%s}", username)

	userID := request.User.UserID
	status := request.User.Status
	item, err := handler.service.Register(userID, status)
	message := ""

	if item.Status == "" {
		message = "User already exists"
	} else {
		message = "Success"
	}

	user := &pb.User{
		UserID: userID, Status: status,
	}
	actionResult := &pb.RegisterResponse{User: user, Message: message}

	return actionResult, err
}

func (handler *ConnectionHandler) CreateConnection(ctx context.Context, request *pb.NewConnectionRequest) (response *pb.NewConnectionResponse, err error) {
	fmt.Println("[ConnectionHandler]:CreateConnection")
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("POST rr: CC {%s}", username)
	return handler.service.CreateConnection(request.Connection.BaseUserUUID, request.Connection.ConnectUserUUID)
}

func (handler *ConnectionHandler) AcceptConnection(ctx context.Context, request *pb.AcceptConnectionMessage) (response *pb.NewConnectionResponse, err error) {
	fmt.Println("[ConnectionHandler]:AcceptConnection")
	return handler.service.AcceptConnection(request.AcceptConnection.RequestSenderUser, request.AcceptConnection.RequestApprovalUser)
}

func (handler *ConnectionHandler) GetAllConnectionForUser(ctx context.Context, request *pb.GetConnectionRequest) (*pb.GetAllResponse, error) {

	users, err := handler.service.GetAllConnectionForUser(request.GetUuid())
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := pb.User{UserID: user.UserUID, Status: string(user.Status)}
		response.Users = append(response.Users, &current)
	}

	return response, nil
}

func (handler *ConnectionHandler) GetAllConnectionRequestsForUser(ctx context.Context,
	request *pb.GetConnectionRequest) (*pb.GetAllResponse, error) {

	users, err := handler.service.GetAllConnectionRequestsForUser(request.GetUuid())
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := pb.User{UserID: user.UserUID, Status: string(user.Status)}
		response.Users = append(response.Users, &current)
	}

	return response, nil
}
