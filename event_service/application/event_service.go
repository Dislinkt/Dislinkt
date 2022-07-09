package application

import (
	"github.com/dislinkt/event_service/domain"
)

type EventService struct {
	store domain.EventStore
}

func NewEventService(store domain.EventStore) *EventService {
	return &EventService{store: store}
}

func (service *EventService) GetAllEvents() ([]*domain.Event, error) {
	return service.store.GetAllEvents()
}

func (service *EventService) InsertEvent(event *domain.Event) error {
	return service.store.InsertEvent(event)
}
