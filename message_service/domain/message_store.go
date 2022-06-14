package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageStore interface {
	GetMessageHistoriesByUser(userId string) ([]*MessageHistory, error) //TODO: get userId from jwt token
	GetMessageHistory(user1Id, user2Id string) (*MessageHistory, error) //TODO: get user1Id from jwt token
	InsertMessage(message *Message, messageHistory string) (*MessageHistory, error)
	GetHistoryById(id primitive.ObjectID) (*MessageHistory, error)
}
