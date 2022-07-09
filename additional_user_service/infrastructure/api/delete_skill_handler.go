package api

import (
	"fmt"
	"github.com/dislinkt/additional_user_service/application"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type DeleteSkillCommandHandler struct {
	additionalService *application.AdditionalUserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewDeleteSkillCommandHandler(additionalService *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*DeleteSkillCommandHandler, error) {
	o := &DeleteSkillCommandHandler{
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

func (handler *DeleteSkillCommandHandler) handle(command *events.DeleteSkillCommand) {
	reply := events.DeleteSkillReply{
		Skill:  command.Skill,
		UserId: command.UserId,
	}

	switch command.Type {
	case events.DeleteSkillInAdditional:
		fmt.Println("additional handler add skill")
		skill, _ := handler.additionalService.FindUserSkill(command.Skill.Id, command.UserId)
		_, err := handler.additionalService.DeleteUserSkill(command.UserId, command.Skill.Id)
		if err != nil {
			fmt.Println("additional handler error not added skill")
			reply.Type = events.AdditionalServiceSkillNotDeleted
			return
		}

		reply.Skill.Name = skill.Name
		reply.Type = events.AdditionalServiceSkillDeleted
		fmt.Println("additional handler add success skill")
		fmt.Println(reply.Skill.Name)
		// reply.Type = events.RegistrationApproved
	case events.RollbackSkillDeleteInAdditional:
		fmt.Println("additional handler-rollback education")
		skill, _ := handler.additionalService.FindUserSkill(command.Skill.Id, command.UserId)
		err, _ := handler.additionalService.CreateSkill(command.UserId, skill)
		if err != nil {
			return
		}
		reply.Type = events.AdditionalSkillDeleteRolledBack
	default:
		reply.Type = events.UnknownDeletedSkillReply
	}

	if reply.Type != events.UnknownDeletedSkillReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
