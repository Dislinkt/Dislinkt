package application

import (
	"fmt"
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type UpdateSkillOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewUpdateSkillOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateSkillOrchestrator,
	error) {
	o := &UpdateSkillOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *UpdateSkillOrchestrator) Start(userId string, skillId string, skill *domain.Skill) error {
	event := &events.UpdateSkillCommand{
		Type: events.UpdateSkillInAdditional,
		Skill: events.SkillUpdate{
			Id:   skillId,
			Name: skill.Name,
		},
		UserId:  userId,
		OldName: "",
	}
	return o.commandPublisher.Publish(event)
}

func (o *UpdateSkillOrchestrator) handle(reply *events.UpdateSkillReply) {
	command := events.UpdateSkillCommand{Skill: reply.Skill, UserId: reply.UserId, OldName: reply.OldName}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownUpdateSkillCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *UpdateSkillOrchestrator) nextCommandType(reply events.UpdateSkillReplyType) events.UpdateSkillCommandType {
	switch reply {
	case events.AdditionalServiceSkillUpdated:
		fmt.Println("proslo u additional skill update")
		return events.UpdateSkillInGraph
	case events.GraphDatabaseSkillNotUpdated:
		return events.RollbackSkillUpdateInAdditional
	case events.AdditionalSkillUpdateRolledBack:
		return events.CancelSkillUpdate
	case events.AdditionalServiceSkillNotUpdated:
		return events.CancelSkillUpdate
	case events.GraphDatabaseSkillUpdated:
		fmt.Println("proslo u graf skill update")
		return events.ApproveSkillUpdating
	default:
		return events.UnknownUpdateSkillCommand
	}
}
