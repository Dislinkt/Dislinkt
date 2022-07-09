package persistence

import (
	"context"
	"github.com/dislinkt/event_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "event"
	COLLECTION = "event"
)

type EventMongoDBStore struct {
	events *mongo.Collection
}

func NewEventMongoDBStore(client *mongo.Client) domain.EventStore {
	events := client.Database(DATABASE).Collection(COLLECTION)
	return &EventMongoDBStore{events: events}
}

func (store *EventMongoDBStore) GetAllEvents() ([]*domain.Event, error) {
	filter := bson.D{}
	return store.filter(filter)
}

func (store *EventMongoDBStore) InsertEvent(event *domain.Event) error {
	_, err := store.events.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}

	return nil
}

func (store *EventMongoDBStore) filter(filter interface{}) ([]*domain.Event, error) {
	cursor, err := store.events.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	return decode(cursor)
}

func (store *EventMongoDBStore) filterOne(filter interface{}) (event *domain.Event, err error) {
	result := store.events.FindOne(context.TODO(), filter)
	err = result.Decode(&event)
	return
}

func decode(cursor *mongo.Cursor) (events []*domain.Event, err error) {
	for cursor.Next(context.TODO()) {
		var event domain.Event
		err = cursor.Decode(&event)
		if err != nil {
			return
		}
		events = append(events, &event)
	}
	err = cursor.Err()
	return
}
