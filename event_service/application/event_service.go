package application

import (
	"context"
	"github.com/dislinkt/common/tracer"
	"github.com/dislinkt/event_service/domain"
)

type EventService struct {
	store domain.EventStore
}

func NewEventService(store domain.EventStore) *EventService {
	return &EventService{store: store}
}

func (service *EventService) GetAllEvents(ctx context.Context) ([]*domain.Event, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllEvents-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAllEvents()
}

func (service *EventService) InsertEvent(ctx context.Context, event *domain.Event) error {
	span := tracer.StartSpanFromContext(ctx, "InsertEvent-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.InsertEvent(event)
}
