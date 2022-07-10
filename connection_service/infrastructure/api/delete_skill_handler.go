package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/connection_service/application"
)

type DeleteSkillCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewDeleteSkillCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*DeleteSkillCommandHandler, error) {
	o := &DeleteSkillCommandHandler{
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

func (handler *DeleteSkillCommandHandler) handle(command *events.DeleteSkillCommand) {
	reply := events.DeleteSkillReply{
		Skill:  command.Skill,
		UserId: command.UserId,
	}
	fmt.Println(command)
	switch command.Type {
	case events.DeleteSkillInGraph:
		fmt.Println("connectionHandler-delete Skill")
		fmt.Println(command.Skill.Name)
		fmt.Println(command.UserId)
		fmt.Println("podaci")
		_, err := handler.connectionService.DeleteSkillToUser(context.TODO(), command.Skill.Name, command.UserId)
		if err != nil {
			fmt.Println("connectionHandler-error delete skill")
			reply.Type = events.GraphDatabaseSkillNotDeleted
			return
		}
		fmt.Println("connectionHandler- delete skill success")
		reply.Type = events.GraphDatabaseSkillDeleted
	default:
		fmt.Println("connectionHandler-unknown reply")
		reply.Type = events.UnknownDeletedSkillReply
	}

	if reply.Type != events.UnknownDeletedSkillReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
