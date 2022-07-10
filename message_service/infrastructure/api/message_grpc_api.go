package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/interceptor"
	pb "github.com/dislinkt/common/proto/message_service"
	notificationGw "github.com/dislinkt/common/proto/notification_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "GetMessageHistoriesByUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	userResponse, _ := persistence.UserClient("user_service:8000").GetUserByUsername(context.TODO(), &userGw.GetOneByUsernameMessage{Username: username})
	messageHistories, err := handler.service.GetMessageHistoriesByUser(ctx, userResponse.User.Id)
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
	span := tracer.StartSpanFromContext(ctx, "GetMessageHistoryAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	userResponse, _ := persistence.UserClient("user_service:8000").GetUserByUsername(context.TODO(), &userGw.GetOneByUsernameMessage{Username: username})
	messageHistory, err := handler.service.GetMessageHistory(ctx, userResponse.User.Id, request.ReceiverId)
	if err != nil {
		return nil, err
	}
	messageHistoryPb := mapMessageHistory(messageHistory)
	response := &pb.GetResponse{MessageHistory: messageHistoryPb}
	return response, nil
}

func (handler *MessageHandler) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (*pb.GetResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "SendMessageAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	message := mapNewMessage(request.Message)
	messageHistory, err := handler.service.InsertMessage(ctx, message, request.MessageHistoryId)
	if err != nil {
		return nil, err
	}
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	_, _ = persistence.NotificationClient("notification_service:8000").SaveNotification(context.TODO(),
		&notificationGw.SaveNotificationRequest{Notification: mapNotification(username), UserId: message.ReceiverId})
	return &pb.GetResponse{MessageHistory: mapMessageHistory(messageHistory)}, nil
}
