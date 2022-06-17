package api

import (
	"context"
	"fmt"
	logger "github.com/dislinkt/common/logging"

	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/patch_user"
	"github.com/dislinkt/connection_service/application"
)

type PatchUserCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
	logger            *logger.Logger
}

func NewPatchUserCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*PatchUserCommandHandler, error) {
	logger := logger.InitLogger(context.TODO())
	o := &PatchUserCommandHandler{
		connectionService: connectionService,
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

func (handler *PatchUserCommandHandler) handle(command *events.PatchUserCommand) {
	reply := events.PatchUserReply{User: command.User}

	switch command.Type {
	case events.PatchUserInConnection:
		fmt.Println("connections handler-patch")
		fmt.Println(command.User)
		// id, _ := uuid.FromString(command.User.Id)
		err := handler.connectionService.UpdateUser(command.User.Id, command.User.Private)
		if err != nil {
			reply.Type = events.PatchFailedInConnection
			return
		}
		reply.Type = events.PatchedInConnection
		handler.logger.InfoLogger.Infof("SC-PU {%s}", reply.User.Username)

	default:
		reply.Type = events.UnknownPatchReply
	}

	if reply.Type != events.UnknownPatchReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
