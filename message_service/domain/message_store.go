package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageStore interface {
	GetMessageHistoriesByUser(userId string) ([]*MessageHistory, error)
	GetMessageHistory(user1Id, user2Id string) (*MessageHistory, error)
	InsertMessage(message *Message, messageHistory string) (*MessageHistory, error)
	GetHistoryById(id primitive.ObjectID) (*MessageHistory, error)
}
