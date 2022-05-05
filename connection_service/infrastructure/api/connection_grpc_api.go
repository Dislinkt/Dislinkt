package api

import (
	"context"
	"fmt"
	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/application"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	service *application.ConnectionService
}

func NewConnectionHandler(service *application.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		service: service,
	}
}

func (handler *ConnectionHandler) Register(ctx context.Context, request *pb.RegisterRequest) (response *pb.RegisterResponse, err error) {
	fmt.Println("[ConnectionHandler]:Register")
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
	return handler.service.CreateConnection(request.Connection.BaseUserUUID, request.Connection.ConnectUserUUID)
}
