package api

import (
	"fmt"
	"github.com/dislinkt/additional_user_service/application"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type AddSkillCommandHandler struct {
	additionalService *application.AdditionalUserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewAddSkillCommandHandler(additionalService *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*AddSkillCommandHandler, error) {
	o := &AddSkillCommandHandler{
		additionalService: additionalService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *AddSkillCommandHandler) handle(command *events.AddSkillCommand) {
	reply := events.AddSkillReply{
		Skill:  command.Skill,
		UserId: command.UserId,
	}

	switch command.Type {
	case events.AddSkillInAdditional:
		fmt.Println("additional handler add skill")
		skill := mapAdditionalCommandAddSkill(command)
		_, err := handler.additionalService.CreateSkill(command.UserId, skill)
		if err != nil {
			fmt.Println("additional handler error not added skill")
			reply.Type = events.AdditionalServiceSkillNotAdded
			return
		}
		reply.UserId = command.UserId
		reply.Type = events.AdditionalServiceSkillAdded
		fmt.Println("additional handler add success skill")
		// reply.Type = events.RegistrationApproved
	case events.RollbackSkillInAdditional:
		fmt.Println("additional handler-rollback education")
		edu := mapAdditionalCommandAddSkill(command)
		err, _ := handler.additionalService.DeleteUserEducation(command.UserId, edu.Id.Hex())
		if err != nil {
			return
		}
		reply.Type = events.AdditionalSkillRolledBack
	default:
		reply.Type = events.UnknownAddSkillReply
	}

	if reply.Type != events.UnknownAddSkillReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
