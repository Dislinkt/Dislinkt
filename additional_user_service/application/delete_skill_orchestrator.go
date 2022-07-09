package application

import (
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type DeleteSkillOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewDeleteSkillOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*DeleteSkillOrchestrator,
	error) {
	o := &DeleteSkillOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *DeleteSkillOrchestrator) Start(uuid string, additionID string) error {
	event := &events.DeleteSkillCommand{
		Type: events.DeleteSkillInAdditional,
		Skill: events.SkillDelete{
			Id:   additionID,
			Name: "",
		},
		UserId: uuid,
	}
	return o.commandPublisher.Publish(event)
}

func (o *DeleteSkillOrchestrator) handle(reply *events.DeleteSkillReply) {
	command := events.DeleteSkillCommand{Skill: reply.Skill, UserId: reply.UserId}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownDeleteSkillCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *DeleteSkillOrchestrator) nextCommandType(reply events.DeleteSkillReplyType) events.DeleteSkillCommandType {
	switch reply {
	case events.AdditionalServiceSkillDeleted:
		fmt.Println("proslo u additional skill delete")
		return events.DeleteSkillInGraph
	case events.GraphDatabaseSkillNotDeleted:
		return events.RollbackSkillDeleteInAdditional
	case events.AdditionalSkillDeleteRolledBack:
		return events.CancelSkillDelete
	case events.AdditionalServiceSkillNotDeleted:
		return events.CancelSkillDelete
	case events.GraphDatabaseSkillDeleted:
		fmt.Println("proslo u graf skill delete")
		return events.ApproveSkillDeleting
	default:
		return events.UnknownDeleteSkillCommand
	}
}
