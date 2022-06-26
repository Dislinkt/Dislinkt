package persistence

import (
	"context"
	"github.com/dislinkt/message_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "message"
	COLLECTION = "message"
)

type MessageMongoDBStore struct {
	messages *mongo.Collection
}

func NewMessagesMongoDBStore(client *mongo.Client) domain.MessageStore {
	messages := client.Database(DATABASE).Collection(COLLECTION)
	return &MessageMongoDBStore{messages: messages}
}

func (store *MessageMongoDBStore) GetMessageHistory(user1Id, user2Id string) (*domain.MessageHistory, error) {
	messageHistory, err := store.getMessageHistoryByUsers(user1Id, user2Id)

	var messages []domain.Message
	for _, message := range messageHistory.Messages {
		if !message.IsRead {
			message.IsRead = true
		}
		messages = append(messages, message)
	}

	_, err = store.messages.UpdateOne(context.TODO(), bson.M{"_id": messageHistory.Id}, bson.D{
		{"$set", bson.D{{"messages", messages}}},
	})
	if err != nil {
		return nil, err
	}

	return messageHistory, err
}

func (store *MessageMongoDBStore) getMessageHistoryByUsers(user1Id string, user2Id string) (*domain.MessageHistory, error) {
	filter := bson.M{
		"$or": []bson.M{
			{
				"user_one_id": user1Id,
				"user_two_id": user2Id,
			},
			{
				"user_one_id": user2Id,
				"user_two_id": user1Id,
			},
		}}

	return store.filterOne(filter)
}

func (store *MessageMongoDBStore) InsertMessage(message *domain.Message, historyId string) (*domain.MessageHistory, error) {
	id, _ := primitive.ObjectIDFromHex(historyId)
	messageHistory, err := store.GetHistoryById(id)

	message.IsRead = false

	if messageHistory == nil {
		history, _ := store.getMessageHistoryByUsers(message.SenderId, message.ReceiverId)
		if history == nil {
			newId := primitive.NewObjectID()
			messHistory := &domain.MessageHistory{
				Id:        newId,
				UserOneId: message.SenderId,
				UserTwoId: message.ReceiverId,
				Messages:  nil,
			}
			messHistory.Messages = append(messHistory.Messages, *message)
			_, err = store.messages.InsertOne(context.TODO(), messHistory)
			if err != nil {
				return nil, err
			}
			messageHistory, err = store.GetHistoryById(newId)
		}
	} else {
		if areUsersMatching(messageHistory, message) {
			messages := append(messageHistory.Messages, *message)

			_, err = store.messages.UpdateOne(context.TODO(), bson.M{"_id": messageHistory.Id}, bson.D{
				{"$set", bson.D{{"messages", messages}}},
			})
			if err != nil {
				return nil, err
			}
			messageHistory, err = store.GetHistoryById(id)
		} else {
			return nil, err
		}
	}

	return messageHistory, err
}

func areUsersMatching(history *domain.MessageHistory, message *domain.Message) bool {
	return (message.SenderId == history.UserOneId && message.ReceiverId == history.UserTwoId) ||
		(message.SenderId == history.UserTwoId && message.ReceiverId == history.UserOneId)
}

func (store *MessageMongoDBStore) GetMessageHistoriesByUser(userId string) ([]*domain.MessageHistory, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"user_one_id": userId},
			{"user_two_id": userId},
		}}

	return store.filter(filter)
}

func (store *MessageMongoDBStore) GetHistoryById(id primitive.ObjectID) (*domain.MessageHistory, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *MessageMongoDBStore) filter(filter interface{}) ([]*domain.MessageHistory, error) {
	cursor, err := store.messages.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	return decode(cursor)
}

func (store *MessageMongoDBStore) filterOne(filter interface{}) (messageHistory *domain.MessageHistory, err error) {
	result := store.messages.FindOne(context.TODO(), filter)
	err = result.Decode(&messageHistory)
	return
}

func decode(cursor *mongo.Cursor) (messageHistories []*domain.MessageHistory, err error) {
	for cursor.Next(context.TODO()) {
		var messageHistory domain.MessageHistory
		err = cursor.Decode(&messageHistory)
		if err != nil {
			return
		}
		messageHistories = append(messageHistories, &messageHistory)
	}
	err = cursor.Err()
	return
}
