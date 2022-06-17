package api

import (
	"context"
	"fmt"
	logger "github.com/dislinkt/common/logging"

	"github.com/dislinkt/auth_service/application"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type UpdateUserCommandHandler struct {
	userService       *application.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
	logger            *logger.Logger
}

func NewUpdateUserCommandHandler(userService *application.UserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*UpdateUserCommandHandler, error) {
	logger := logger.InitLogger(context.TODO())
	o := &UpdateUserCommandHandler{
		userService:       userService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
		logger:            logger,
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
	case events.UpdateInAuth:
		fmt.Println("auth handler-update")

		err := handler.userService.ChangeUsername(command.User.Id, command.User.Username)
		if err != nil {
			// reply.User =
			reply.Type = events.UserNotUpdatedInAuth
			return
		}
		reply.Type = events.UserUpdatedInAuth
		handler.logger.InfoLogger.Infof("SC-UU {%s}", reply.User.Username)

	default:
		reply.Type = events.UnknownUpdateReply
	}

	if reply.Type != events.UnknownUpdateReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
