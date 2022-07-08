package api

import (
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/connection_service/application"
)

type CreateJobOfferCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewJobOfferCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateJobOfferCommandHandler, error) {
	o := &CreateJobOfferCommandHandler{
		connectionService: connectionService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateJobOfferCommandHandler) handle(command *events.CreateJobOfferCommand) {
	reply := events.CreateJobReply{JobOffer: command.JobOffer}

	switch command.Type {
	case events.CreateJobOfferInGraph:
		fmt.Println("connectionHandler-createJobOffer")
		jobOffer := mapConnectionCommandCreateJob(command)

		_, err := handler.connectionService.InsertJobOffer(*jobOffer)
		if err != nil {
			fmt.Println("connectionHandler-error")
			reply.Type = events.ConnectionServiceNotCreated
			return
		}
		fmt.Println("connectionHandler-createJobOffer1")
		reply.Type = events.ConnectionServiceCreated
	default:
		fmt.Println("connectionHandler-unknown")
		reply.Type = events.UnknownCreateJobOfferReply
	}

	if reply.Type != events.UnknownCreateJobOfferReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
