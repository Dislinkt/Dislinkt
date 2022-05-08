package api

import (
	"fmt"

	"github.com/dislinkt/additional_user_service/application"
	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/register_user"
)

type CreateUserCommandHandler struct {
	additionalService *application.AdditionalUserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(additionalService *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
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

func (handler *CreateUserCommandHandler) handle(command *events.RegisterUserCommand) {
	reply := events.RegisterUserReply{User: command.User}

	switch command.Type {
	case events.UpdateAdditional:
		fmt.Println("additional handler-update")
		err := handler.additionalService.CreateDocument(command.User.Id)
		if err != nil {
			reply.Type = events.AdditionalServiceNotUpdated
			return
		}
		reply.Type = events.AdditionalServiceUpdated
		//reply.Type = events.RegistrationApproved
	case events.RollbackAdditional:
		fmt.Println("additional handler-rollback")
		err := handler.additionalService.DeleteDocument(command.User.Id)
		if err != nil {
			return
		}
		reply.Type = events.AdditionalServiceRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
