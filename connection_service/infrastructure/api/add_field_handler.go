package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/connection_service/application"
)

type AddEducationCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewAddEducationCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*AddEducationCommandHandler, error) {
	o := &AddEducationCommandHandler{
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

func (handler *AddEducationCommandHandler) handle(command *events.AddEducationCommand) {
	reply := events.AddEducationReply{
		Education: command.Education,
		UserId:    command.UserId,
	}
	fmt.Println(command)
	fmt.Println(reply)
	fmt.Println(command.Type)
	switch command.Type {
	case events.AddEducationInGraph:
		fmt.Println("connectionHandler-addEducation")
		fmt.Println(command.Education.FieldOfStudy)
		fmt.Println(command.UserId)
		fmt.Println("podaci")
		_, err := handler.connectionService.InsertFieldToUser(context.TODO(), command.Education.FieldOfStudy, command.UserId)
		if err != nil {
			fmt.Println("connectionHandler-error addEducation")
			reply.Type = events.GraphDatabaseNotAdded
			return
		}
		fmt.Println("connectionHandler- addEducation success")
		reply.Type = events.GraphDatabaseAdded
	default:
		fmt.Println("connectionHandler-unknown reply")
		reply.Type = events.UnknownAddEducationReply
	}

	if reply.Type != events.UnknownAddEducationReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
