package application

import (
	"context"
	"github.com/dislinkt/common/tracer"
	"github.com/dislinkt/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	store domain.MessageStore
}

func NewMessageService(store domain.MessageStore) *MessageService {
	return &MessageService{store: store}
}

func (service *MessageService) GetMessageHistory(ctx context.Context, user1Id, user2Id string) (*domain.MessageHistory, error) {
	span := tracer.StartSpanFromContext(ctx, "GetMessageHistory-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetMessageHistory(user1Id, user2Id)
}

func (service *MessageService) InsertMessage(ctx context.Context, message *domain.Message, historyId string) (*domain.MessageHistory, error) {
	span := tracer.StartSpanFromContext(ctx, "InsertMessage-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.InsertMessage(message, historyId)
}

func (service *MessageService) GetMessageHistoriesByUser(ctx context.Context, userId string) ([]*domain.MessageHistory, error) {
	span := tracer.StartSpanFromContext(ctx, "GetMessageHistoriesByUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetMessageHistoriesByUser(userId)
}

func (service *MessageService) GetHistoryById(ctx context.Context, id primitive.ObjectID) (*domain.MessageHistory, error) {
	span := tracer.StartSpanFromContext(ctx, "GetHistoryById")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetHistoryById(id)
}
