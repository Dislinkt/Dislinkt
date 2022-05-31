package api

import (
	"fmt"

	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/user_service/application"
)

type UpdateUserCommandHandler struct {
	userService       *application.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateUserCommandHandler(userService *application.UserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*UpdateUserCommandHandler, error) {
	o := &UpdateUserCommandHandler{
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

func (handler *UpdateUserCommandHandler) handle(command *events.UpdateUserCommand) {
	reply := events.UpdateUserReply{User: command.User}

	switch command.Type {
	case events.UpdateInUser:
		fmt.Println("update user handler-update")
		fmt.Println(command.User)
		user := mapCommandUpdateUser(command)
		_, err := handler.userService.Update(user.Id, user)
		if err != nil {
			reply.Type = events.UserNotUpdatedInUser
			return
		}
		reply.Type = events.UserUpdatedInUser
	case events.RollbackUpdateInUser:
		fmt.Println("update user handler-rollback")
		user := mapCommandUpdateUser(command)
		_, err := handler.userService.Update(user.Id, user)
		if err != nil {
			return
		}
		reply.Type = events.UserRolledBackInUser
	default:
		reply.Type = events.UnknownUpdateReply
	}

	if reply.Type != events.UnknownUpdateReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
