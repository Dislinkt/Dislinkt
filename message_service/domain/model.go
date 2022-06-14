package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageHistory struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserOneId string             `bson:"user_one_id"`
	UserTwoId string             `bson:"user_two_id"`
	Messages  []Message          `bson:"messages"`
}

type Message struct {
	SenderId    string    `bson:"sender_one_id"`
	ReceiverId  string    `bson:"receiver_two_id"`
	MessageText string    `bson:"receiver_two_id"`
	DateSent    time.Time `bson:"date_posted"`
}
