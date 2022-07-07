package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	Id               primitive.ObjectID `bson:"_id"`
	NotificationText string             `bson:"notification_text"`
	NotificationType NotificationType   `bson:"notification_type"`
	UserId           string             `bson:"user_id"`
	Date             time.Time          `bson:"date"`
	IsRead           bool               `bson:"is_read"`
}

type NotificationType int

const (
	Unknown NotificationType = iota
	CONNECTION
	MESSAGE
	POST
)
