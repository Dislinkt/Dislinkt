package application

import (
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/domain"
)

type CreateJobOfferOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreateJobOfferOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*CreateJobOfferOrchestrator,
	error) {
	o := &CreateJobOfferOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *CreateJobOfferOrchestrator) Start(jobOffer *domain.JobOffer) error {
	event := &events.CreateJobOfferCommand{
		Type: events.CreateJobOfferInPost,
		JobOffer: events.JobOffer{
			Id:            primitive.NewObjectID(),
			Position:      jobOffer.Position,
			Description:   jobOffer.Description,
			Preconditions: jobOffer.Preconditions,
			DatePosted:    jobOffer.DatePosted,
			Duration:      jobOffer.Duration,
			Location:      jobOffer.Location,
			Title:         jobOffer.Title,
			Field:         jobOffer.Field,
		},
	}
	return o.commandPublisher.Publish(event)
}

func (o *CreateJobOfferOrchestrator) handle(reply *events.CreateJobReply) {
	command := events.CreateJobOfferCommand{JobOffer: reply.JobOffer}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCreateJobOfferCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *CreateJobOfferOrchestrator) nextCommandType(reply events.CreateJobOfferReplyType) events.CreateJobOfferCommandType {
	switch reply {
	case events.PostServiceCreated:
		fmt.Println("uslo orc")
		return events.CreateJobOfferInGraph
	case events.ConnectionServiceNotCreated:
		return events.RollbackJobOfferInPost
	case events.PostServiceRolledBack:
		return events.CancelJobOfferCreation
	case events.PostServiceNotCreated:
		return events.CancelJobOfferCreation
	case events.ConnectionServiceCreated:
		fmt.Println("uslo orc1")
		return events.ApproveCreation
	default:
		return events.UnknownCreateJobOfferCommand
	}
}
