package application

import (
	"fmt"
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type AddSkillOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewAddSkillOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*AddSkillOrchestrator,
	error) {
	o := &AddSkillOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *AddSkillOrchestrator) Start(skill *domain.Skill, userId string) error {
	event := &events.AddSkillCommand{
		Type: events.AddSkillInAdditional,
		Skill: events.Skill{
			Id:   skill.Id,
			Name: skill.Name,
		},
		UserId: userId,
	}
	return o.commandPublisher.Publish(event)
}

func (o *AddSkillOrchestrator) handle(reply *events.AddSkillReply) {
	command := events.AddSkillCommand{Skill: reply.Skill, UserId: reply.UserId}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownAddSkillCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *AddSkillOrchestrator) nextCommandType(reply events.AddSkillReplyType) events.AddSkillCommandType {
	switch reply {
	case events.AdditionalServiceSkillAdded:
		fmt.Println("proslo u additional skill")
		return events.AddSkillInGraph
	case events.GraphDatabaseSkillNotAdded:
		return events.RollbackSkillInAdditional
	case events.AdditionalSkillRolledBack:
		return events.CancelSkillAdd
	case events.AdditionalServiceSkillNotAdded:
		return events.CancelSkillAdd
	case events.GraphDatabaseSkillAdded:
		fmt.Println("proslo u graf skill")
		return events.ApproveSkillAdding
	default:
		return events.UnknownAddSkillCommand
	}
}
