package api

import (
	"github.com/dislinkt/auth_service/application"
	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/register_user"
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
	case events.UpdateAuth:
		uuid, err := handler.userService.Insert(mapCommandUser(command))
		reply.User.Id = uuid.String()
		if err != nil {
			reply.Type = events.AuthNotUpdated
			return
		}
		reply.Type = events.AuthUpdated
	case events.RollbackAuth:
		err := handler.userService.Delete(mapCommandUser(command))
		if err != nil {
			return
		}
		reply.Type = events.AuthRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
