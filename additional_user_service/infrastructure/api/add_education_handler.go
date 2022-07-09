package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/additional_user_service/application"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type AddEducationCommandHandler struct {
	additionalService *application.AdditionalUserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewAddEducationCommandHandler(additionalService *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*AddEducationCommandHandler, error) {
	o := &AddEducationCommandHandler{
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

func (handler *AddEducationCommandHandler) handle(command *events.AddEducationCommand) {
	reply := events.AddEducationReply{
		Education: command.Education,
		UserId:    command.UserId,
	}

	switch command.Type {
	case events.AddEducationInAdditional:
		fmt.Println("additional handler add education")
		edu := mapAdditionalCommandAddEducation(command)
		_, err := handler.additionalService.CreateEducation(context.TODO(), command.UserId, edu)
		if err != nil {
			fmt.Println("additional handler error not added")
			reply.Type = events.AdditionalServiceNotAdded
			return
		}
		reply.UserId = command.UserId
		reply.Type = events.AdditionalServiceAdded
		fmt.Println("additional handler add success")
		// reply.Type = events.RegistrationApproved
	case events.RollbackEducationInAdditional:
		fmt.Println("additional handler-rollback education")
		edu := mapAdditionalCommandAddEducation(command)
		err, _ := handler.additionalService.DeleteUserEducation(context.TODO(), command.UserId, edu.Id.Hex())
		if err != nil {
			return
		}
		reply.Type = events.AdditionalRolledBack
	default:
		reply.Type = events.UnknownAddEducationReply
	}

	if reply.Type != events.UnknownAddEducationReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
