package api

import (
	pb "github.com/dislinkt/common/proto/event_service"
	"github.com/dislinkt/event_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func mapEvent(event *domain.Event) *pb.Event {
	eventPb := &pb.Event{
		UserId:      event.UserId,
		Description: event.Description,
		Date:        event.Date.String(),
	}

	return eventPb
}

func mapNewEvent(eventPb *pb.NewEvent) *domain.Event {
	event := &domain.Event{
		Id:          primitive.NewObjectID(),
		UserId:      eventPb.UserId,
		Description: eventPb.Description,
		Date:        time.Now(),
	}
	return event
}
