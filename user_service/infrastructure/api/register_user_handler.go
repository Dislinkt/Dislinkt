package api

import (
	"fmt"

	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/register_user"
	"github.com/dislinkt/user_service/application"
)

type CreateUserCommandHandler struct {
	userService       *application.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(userService *application.UserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		userService:       userService,
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
	case events.UpdateUser:
		fmt.Println("userhandler-update")
		err := handler.userService.Insert(mapCommandUser(command))
		if err != nil {
			reply.Type = events.UserServiceNotUpdated
			return
		}
		reply.Type = events.UserServiceUpdated
	case events.RollbackUser:
		fmt.Println("userhandler-rollback")
		err := handler.userService.Delete(mapCommandUser(command))
		if err != nil {
			return
		}
		reply.Type = events.UserServiceRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
