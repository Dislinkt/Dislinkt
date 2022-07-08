package application

import (
	"fmt"
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddEducationOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewAddEducationOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*AddEducationOrchestrator,
	error) {
	o := &AddEducationOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *AddEducationOrchestrator) Start(education *domain.Education, userId string) error {
	event := &events.AddEducationCommand{
		Type: events.AddEducationInAdditional,
		Education: events.Education{
			Id:           primitive.NewObjectID(),
			School:       education.School,
			Degree:       education.Degree,
			FieldOfStudy: education.FieldOfStudy,
			StartDate:    education.StartDate,
			EndDate:      education.EndDate,
		},
		UserId: userId,
	}
	return o.commandPublisher.Publish(event)
}

func (o *AddEducationOrchestrator) handle(reply *events.AddEducationReply) {
	command := events.AddEducationCommand{Education: reply.Education, UserId: reply.UserId}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownAddEducationCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *AddEducationOrchestrator) nextCommandType(reply events.AddEducationReplyType) events.AddEducationCommandType {
	switch reply {
	case events.AdditionalServiceAdded:
		fmt.Println("proslo u additional")
		return events.AddEducationInGraph
	case events.GraphDatabaseNotAdded:
		return events.RollbackEducationInAdditional
	case events.AdditionalRolledBack:
		return events.CancelEducationAdd
	case events.AdditionalServiceNotAdded:
		return events.CancelEducationAdd
	case events.GraphDatabaseAdded:
		fmt.Println("proslo u graf")
		return events.ApproveAdding
	default:
		return events.UnknownAddEducationCommand
	}
}
