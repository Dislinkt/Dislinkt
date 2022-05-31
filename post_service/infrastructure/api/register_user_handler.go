package api

import (
	"fmt"

	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/register_user"
	"post_service/application"
)

type CreateUserCommandHandler struct {
	postService       *application.PostService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(postService *application.PostService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		postService:       postService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateUserCommandHandler) handle(command *events.RegisterUserCommand) {
	reply := events.RegisterUserReply{User: command.User}

	switch command.Type {
	case events.UpdatePost:
		fmt.Println("posthandler-update")

		// Napravila sam ti saga_mapper samo prilagodi kako tebi odgovara

		// err := handler.postService.InsertUser(mapCommandUser(command))
		// if err != nil {
		// 	reply.Type = events.PostNotUpdated
		// 	return
		// }
		reply.Type = events.PostUpdated
	case events.RollbackPost:
		fmt.Println("posthandler-rollback")
		// err := handler.postService.DeleteUser(mapCommandUser(command))
		// if err != nil {
		// 	return
		// }
		reply.Type = events.PostRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
