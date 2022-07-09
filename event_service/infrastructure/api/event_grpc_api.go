package api

import (
	"context"
	pb "github.com/dislinkt/common/proto/event_service"
	"github.com/dislinkt/event_service/application"
)

type EventHandler struct {
	pb.UnimplementedEventServiceServer
	service *application.EventService
}

func NewEventHandler(service *application.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (handler *EventHandler) GetAllEvents(ctx context.Context, request *pb.Empty) (*pb.GetMultipleResponse, error) {
	events, err := handler.service.GetAllEvents()
	if err != nil {
		return nil, err
	}
	response := &pb.GetMultipleResponse{Events: []*pb.Event{}}
	for _, event := range events {
		current := mapEvent(event)
		response.Events = append(response.Events, current)
	}

	for i, j := 0, len(response.Events)-1; i < j; i, j = i+1, j-1 {
		response.Events[i], response.Events[j] = response.Events[j], response.Events[i]
	}

	return response, nil
}

func (handler *EventHandler) SaveEvent(ctx context.Context, request *pb.SaveEventRequest) (*pb.Empty, error) {
	event := mapNewEvent(request.Event)
	err := handler.service.InsertEvent(event)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
