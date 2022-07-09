package api

import (
	"context"
	"fmt"
	eventGw "github.com/dislinkt/common/proto/event_service"
	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/patch_user"
	"github.com/dislinkt/connection_service/application"
	"github.com/dislinkt/connection_service/infrastructure/persistance"
)

type PatchUserCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewPatchUserCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*PatchUserCommandHandler, error) {
	o := &PatchUserCommandHandler{
		connectionService: connectionService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
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

		_, _ = persistance.EventClient("event_service:8000").SaveEvent(context.TODO(),
			&eventGw.SaveEventRequest{Event: mapEventForUserPrivacyChange(command.User.Id, command.User.Private)})

		reply.Type = events.PatchedInConnection

	default:
		reply.Type = events.UnknownPatchReply
	}

	if reply.Type != events.UnknownPatchReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
