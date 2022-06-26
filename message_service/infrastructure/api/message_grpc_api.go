package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/interceptor"
	pb "github.com/dislinkt/common/proto/message_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/message_service/application"
	"github.com/dislinkt/message_service/infrastructure/persistence"
)

type MessageHandler struct {
	pb.UnimplementedMessageServiceServer
	service *application.MessageService
}

func NewMessageHandler(service *application.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (handler *MessageHandler) GetMessageHistoriesByUser(ctx context.Context, request *pb.Empty) (*pb.GetMultipleResponse, error) {
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	userResponse, _ := persistence.UserClient("user_service:8000").GetUserByUsername(context.TODO(), &userGw.GetOneByUsernameMessage{Username: username})
	messageHistories, err := handler.service.GetMessageHistoriesByUser(userResponse.User.Id)
	if err != nil {
		return nil, err
	}
	response := &pb.GetMultipleResponse{MessageHistories: []*pb.MessageHistory{}}
	for _, messHistory := range messageHistories {
		current := mapMessageHistory(messHistory)
		response.MessageHistories = append(response.MessageHistories, current)
	}
	return response, nil
}

func (handler *MessageHandler) GetMessageHistory(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	userResponse, _ := persistence.UserClient("user_service:8000").GetUserByUsername(context.TODO(), &userGw.GetOneByUsernameMessage{Username: username})
	messageHistory, err := handler.service.GetMessageHistory(userResponse.User.Id, request.ReceiverId)
	if err != nil {
		return nil, err
	}
	messageHistoryPb := mapMessageHistory(messageHistory)
	response := &pb.GetResponse{MessageHistory: messageHistoryPb}
	return response, nil
}

func (handler *MessageHandler) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (*pb.GetResponse, error) {
	message := mapNewMessage(request.Message)
	messageHistory, err := handler.service.InsertMessage(message, request.MessageHistoryId)
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{MessageHistory: mapMessageHistory(messageHistory)}, nil
}
