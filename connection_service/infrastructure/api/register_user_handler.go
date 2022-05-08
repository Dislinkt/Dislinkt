package api

import (
	"fmt"

	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/register_user"
	"github.com/dislinkt/connection_service/application"
)

type CreateUserCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(userService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		connectionService: userService,
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

	var status string
	if command.User.Private {
		status = "PRIVATE"
	} else {
		status = "PUBLIC"
	}

	switch command.Type {
	case events.UpdateConnectionNode:
		fmt.Println("connection handler-update")
		_, err := handler.connectionService.Register(command.User.Id, status)
		if err != nil {
			reply.Type = events.ConnectionsNotUpdated
			return
		}
		reply.Type = events.ConnectionsUpdated
	case events.RollbackConnectionNode:
		fmt.Println("connection handler-rollback")
		// err := handler.connectionService.Delete(mapCommandUser(command))
		// if err != nil {
		// 	return
		// }
		reply.Type = events.UserServiceRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
