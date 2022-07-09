package api

import (
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/connection_service/application"
)

type AddSkillCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewAddSkillCommandHandler(connectionService *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*AddSkillCommandHandler, error) {
	o := &AddSkillCommandHandler{
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

func (handler *AddSkillCommandHandler) handle(command *events.AddSkillCommand) {
	reply := events.AddSkillReply{
		Skill:  command.Skill,
		UserId: command.UserId,
	}
	fmt.Println(command)
	fmt.Println(reply)
	fmt.Println(command.Type)
	switch command.Type {
	case events.AddSkillInGraph:
		fmt.Println("connectionHandler-add Skill")
		fmt.Println(command.Skill.Name)
		fmt.Println(command.UserId)
		fmt.Println("podaci")
		_, err := handler.connectionService.InsertSkillToUser(command.Skill.Name, command.UserId)
		if err != nil {
			fmt.Println("connectionHandler-error add skill")
			reply.Type = events.GraphDatabaseSkillNotAdded
			return
		}
		fmt.Println("connectionHandler- add skill success")
		reply.Type = events.GraphDatabaseSkillAdded
	default:
		fmt.Println("connectionHandler-unknown reply")
		reply.Type = events.UnknownAddSkillReply
	}

	if reply.Type != events.UnknownAddSkillReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
