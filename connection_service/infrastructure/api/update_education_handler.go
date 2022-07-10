package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/connection_service/application"
)

type UpdateEducationCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateEducationCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*UpdateEducationCommandHandler, error) {
	o := &UpdateEducationCommandHandler{
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

func (handler *UpdateEducationCommandHandler) handle(command *events.UpdateEducationCommand) {
	reply := events.UpdateEducationReply{
		Education:    command.Education,
		UserId:       command.UserId,
		OldFieldName: command.OldFieldName,
	}
	fmt.Println(command)
	switch command.Type {
	case events.UpdateEducationInGraph:
		fmt.Println("connectionHandler-update Education")
		fmt.Println(command.Education.FieldOfStudy)
		fmt.Println(command.UserId)
		fmt.Println(command.OldFieldName)
		fmt.Println("podaci")
		_, err := handler.connectionService.UpdateEducationForUser(context.TODO(), command.UserId, command.OldFieldName, command.Education.FieldOfStudy)
		if err != nil {
			fmt.Println("connectionHandler-error update Education")
			reply.Type = events.GraphDatabaseEducationNotUpdated
			return
		}
		fmt.Println("connectionHandler- update Education success")
		reply.Type = events.GraphDatabaseEducationUpdated
	default:
		fmt.Println("connectionHandler-unknown reply")
		reply.Type = events.UnknownUpdatedEducationReply
	}

	if reply.Type != events.UnknownUpdatedEducationReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
