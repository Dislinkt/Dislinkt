package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/connection_service/application"
)

type DeleteEducationCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewDeleteEducationCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*DeleteEducationCommandHandler, error) {
	o := &DeleteEducationCommandHandler{
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

func (handler *DeleteEducationCommandHandler) handle(command *events.DeleteEducationCommand) {
	reply := events.DeleteEducationReply{
		Education: command.Education,
		UserId:    command.UserId,
	}
	fmt.Println(command)
	switch command.Type {
	case events.DeleteEducationInGraph:
		fmt.Println("connectionHandler-delete Education")
		fmt.Println(command.Education.FieldOfStudy)
		fmt.Println(command.UserId)
		fmt.Println("podaci")
		_, err := handler.connectionService.DeleteFieldToUser(context.TODO(), command.Education.FieldOfStudy, command.UserId)
		if err != nil {
			fmt.Println("connectionHandler-error delete skill")
			reply.Type = events.GraphDatabaseEducationNotDeleted
			return
		}
		fmt.Println("connectionHandler- delete skill success")
		reply.Type = events.GraphDatabaseEducationDeleted
	default:
		fmt.Println("connectionHandler-unknown reply")
		reply.Type = events.UnknownDeletedEducationReply
	}

	if reply.Type != events.UnknownDeletedEducationReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
