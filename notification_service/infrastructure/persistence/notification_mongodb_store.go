package persistence

import (
	"context"
	"github.com/dislinkt/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "notification"
	COLLECTION = "notification"
)

type NotificationMongoDBStore struct {
	notifications *mongo.Collection
}

func NewNotificationMongoDBStore(client *mongo.Client) domain.NotificationStore {
	notifications := client.Database(DATABASE).Collection(COLLECTION)
	return &NotificationMongoDBStore{notifications: notifications}
}

func (store *NotificationMongoDBStore) GetNotificationsForUser(userId string) ([]*domain.Notification, error) {
	filter := bson.M{"user_id": userId}

	notifications, err := store.filter(filter)
	if err != nil {
		return nil, err
	}
	for _, notification := range notifications {
		if !notification.IsRead {
			_, err := store.notifications.UpdateOne(context.TODO(), bson.M{"_id": notification.Id}, bson.D{
				{"$set", bson.D{{"is_read", true}}},
			},
			)
			if err != nil {
				return notifications, err
			}

		}
	}

	return notifications, nil
}

func (store *NotificationMongoDBStore) InsertNotification(notification *domain.Notification) error {
	_, err := store.notifications.InsertOne(context.TODO(), notification)
	if err != nil {
		return err
	}

	return nil
}

func (store *NotificationMongoDBStore) filter(filter interface{}) ([]*domain.Notification, error) {
	cursor, err := store.notifications.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	return decode(cursor)
}

func (store *NotificationMongoDBStore) filterOne(filter interface{}) (notification *domain.Notification, err error) {
	result := store.notifications.FindOne(context.TODO(), filter)
	err = result.Decode(&notification)
	return
}

func decode(cursor *mongo.Cursor) (notifications []*domain.Notification, err error) {
	for cursor.Next(context.TODO()) {
		var notification domain.Notification
		err = cursor.Decode(&notification)
		if err != nil {
			return
		}
		notifications = append(notifications, &notification)
	}
	err = cursor.Err()
	return
}
