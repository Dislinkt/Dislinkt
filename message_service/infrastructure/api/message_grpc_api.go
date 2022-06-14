package api

import (
	"context"
	pb "github.com/dislinkt/common/proto/message_service"
	"github.com/dislinkt/message_service/application"
)

type MessageHandler struct {
	pb.UnimplementedMessageServiceServer
	service *application.MessageService
}

func NewMessageHandler(service *application.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (handler *MessageHandler) GetMessageHistoriesByUser(ctx context.Context, request *pb.Empty) (*pb.GetMultipleResponse, error) {
	userId := "1" //TODO: get from jwt token
	messageHistories, err := handler.service.GetMessageHistoriesByUser(userId)
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
	messageHistory, err := handler.service.GetMessageHistory("1", request.ReceiverId) //TODO: get first id from jwt token
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
