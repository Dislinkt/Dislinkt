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
	filter := bson.M{
		"$or": []bson.M{
			{"$and": []bson.M{
				{"user_one_id": user1Id},
				{"user_two_id": user2Id},
			}},
			{"$and": []bson.M{
				{"user_one_id": user2Id},
				{"user_two_id": user1Id},
			}},
		},
	}

	return store.filterOne(filter)
}

func (store *MessageMongoDBStore) InsertMessage(message *domain.Message, historyId string) (*domain.MessageHistory, error) {
	id, _ := primitive.ObjectIDFromHex(historyId)
	messageHistory, err := store.GetHistoryById(id)

	if messageHistory == nil {
		messHistory := &domain.MessageHistory{
			Id:        primitive.NewObjectID(),
			UserOneId: message.SenderId,
			UserTwoId: message.ReceiverId,
			Messages:  nil,
		}
		result, err := store.messages.InsertOne(context.TODO(), messHistory)
		if err != nil {
			return nil, err
		}

		filter := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
		messageHistory, err = store.filterOne(filter)
	} else {
		messages := append(messageHistory.Messages, *message)

		_, err := store.messages.UpdateOne(context.TODO(), bson.M{"_id": messageHistory.Id}, bson.D{
			{"$set", bson.D{{"messages", messages}}},
		})
		if err != nil {
			return nil, err
		}
	}
	messageHistory.Messages = append(messageHistory.Messages, *message)

	return messageHistory, err
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
