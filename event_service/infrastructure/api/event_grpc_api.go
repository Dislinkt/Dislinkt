package api

import (
	"context"
	pb "github.com/dislinkt/common/proto/event_service"
	"github.com/dislinkt/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "GetAllEventsAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	events, err := handler.service.GetAllEvents(ctx)
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
	span := tracer.StartSpanFromContext(ctx, "SaveEventAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	event := mapNewEvent(request.Event)
	err := handler.service.InsertEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
