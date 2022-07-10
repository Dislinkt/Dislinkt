package application

import (
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type DeleteEducationOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewDeleteEducationOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*DeleteEducationOrchestrator,
	error) {
	o := &DeleteEducationOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *DeleteEducationOrchestrator) Start(uuid string, additionID string) error {
	event := &events.DeleteEducationCommand{
		Type: events.DeleteEducationInAdditional,
		Education: events.EducationDelete{
			Id:           additionID,
			FieldOfStudy: "",
		},
		UserId: uuid,
	}
	return o.commandPublisher.Publish(event)
}

func (o *DeleteEducationOrchestrator) handle(reply *events.DeleteEducationReply) {
	command := events.DeleteEducationCommand{Education: reply.Education, UserId: reply.UserId}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownDeleteEducationCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *DeleteEducationOrchestrator) nextCommandType(reply events.DeleteEducationReplyType) events.DeleteEducationCommandType {
	switch reply {
	case events.AdditionalServiceEducationDeleted:
		fmt.Println("proslo u additional education delete")
		return events.DeleteEducationInGraph
	case events.GraphDatabaseEducationNotDeleted:
		return events.RollbackDeleteEducationInAdditional
	case events.AdditionalEducationDeleteRolledBack:
		return events.CancelEducationDelete
	case events.AdditionalServiceEducationNotDeleted:
		return events.CancelEducationDelete
	case events.GraphDatabaseEducationDeleted:
		fmt.Println("proslo u graf education delete")
		return events.ApproveEducationDeleting
	default:
		return events.UnknownDeleteEducationCommand
	}
}
