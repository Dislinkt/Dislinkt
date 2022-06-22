package application

import (
	"github.com/dislinkt/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	store domain.MessageStore
}

func NewMessageService(store domain.MessageStore) *MessageService {
	return &MessageService{store: store}
}

func (service *MessageService) GetMessageHistory(user1Id, user2Id string) (*domain.MessageHistory, error) {
	return service.store.GetMessageHistory(user1Id, user2Id)
}

func (service *MessageService) InsertMessage(message *domain.Message, historyId string) (*domain.MessageHistory, error) {
	return service.store.InsertMessage(message, historyId)
}

func (service *MessageService) GetMessageHistoriesByUser(userId string) ([]*domain.MessageHistory, error) {
	return service.store.GetMessageHistoriesByUser(userId)
}

func (service *MessageService) GetHistoryById(id primitive.ObjectID) (*domain.MessageHistory, error) {
	return service.store.GetHistoryById(id)
}
