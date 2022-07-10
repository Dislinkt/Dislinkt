package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/connection_service/application"
)

type UpdateSkillCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateSkillCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*UpdateSkillCommandHandler, error) {
	o := &UpdateSkillCommandHandler{
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

func (handler *UpdateSkillCommandHandler) handle(command *events.UpdateSkillCommand) {
	reply := events.UpdateSkillReply{
		Skill:   command.Skill,
		UserId:  command.UserId,
		OldName: command.OldName,
	}
	fmt.Println(command)
	switch command.Type {
	case events.UpdateSkillInGraph:
		fmt.Println("connectionHandler-update Skill")
		fmt.Println(command.Skill.Name)
		fmt.Println(command.UserId)
		fmt.Println(command.OldName)
		fmt.Println("podaci")
		_, err := handler.connectionService.UpdateSkillForUser(context.TODO(), command.UserId, command.OldName, command.Skill.Name)
		if err != nil {
			fmt.Println("connectionHandler-error update skill")
			reply.Type = events.GraphDatabaseSkillNotUpdated
			return
		}
		fmt.Println("connectionHandler- update skill success")
		reply.Type = events.GraphDatabaseSkillUpdated
	default:
		fmt.Println("connectionHandler-unknown reply")
		reply.Type = events.UnknownUpdatedSkillReply
	}

	if reply.Type != events.UnknownUpdatedSkillReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
