package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	Id          primitive.ObjectID `bson:"_id"`
	UserId      string             `bson:"user_id"`
	Description string             `bson:"description"`
	Date        time.Time          `bson:"date"`
}
