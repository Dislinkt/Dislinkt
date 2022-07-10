package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/additional_user_service/application"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type UpdateSkillCommandHandler struct {
	additionalService *application.AdditionalUserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateSkillCommandHandler(additionalService *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*UpdateSkillCommandHandler, error) {
	o := &UpdateSkillCommandHandler{
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

func (handler *UpdateSkillCommandHandler) handle(command *events.UpdateSkillCommand) {
	reply := events.UpdateSkillReply{
		Skill:  command.Skill,
		UserId: command.UserId,
	}

	switch command.Type {
	case events.UpdateSkillInAdditional:
		fmt.Println("additional handler add skill")
		skill, _ := handler.additionalService.FindUserSkill(context.TODO(), command.Skill.Id, command.UserId)
		skill_update := mapAdditionalCommandUpdateSkill(command)
		_, err := handler.additionalService.UpdateUserSkill(context.TODO(), command.UserId, command.Skill.Id, skill_update)
		if err != nil {
			fmt.Println("additional handler error not added skill")
			reply.Type = events.AdditionalServiceSkillNotUpdated
			return
		}

		if skill.Name != skill_update.Name {
			reply.OldName = skill.Name
			reply.Type = events.AdditionalServiceSkillUpdated
		} else {
			reply.Type = events.GraphDatabaseSkillUpdated
		}

		fmt.Println("additional handler update success skill")
		fmt.Println(reply.Skill.Name)
		fmt.Println(reply.OldName)
		// reply.Type = events.RegistrationApproved
	case events.RollbackSkillUpdateInAdditional:
		fmt.Println("additional handler-rollback skill update")
		skill, _ := handler.additionalService.FindUserSkill(context.TODO(), command.Skill.Id, command.UserId)
		err, _ := handler.additionalService.CreateSkill(context.TODO(), command.UserId, skill)
		if err != nil {
			return
		}
		reply.Type = events.AdditionalSkillUpdateRolledBack
	default:
		reply.Type = events.UnknownUpdatedSkillReply
	}

	if reply.Type != events.UnknownUpdatedSkillReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
