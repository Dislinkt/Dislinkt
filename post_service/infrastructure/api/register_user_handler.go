package api

import (
	"context"
	"fmt"

	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
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

		err := handler.postService.InsertUser(context.TODO(), mapPostCommandUser(command))
		if err != nil {
			fmt.Println(err)
			reply.Type = events.PostNotUpdated
			return
		}
		reply.Type = events.PostUpdated
	case events.RollbackPost:
		fmt.Println("posthandler-rollback")
		err := handler.postService.DeleteUser(context.TODO(), mapPostCommandUser(command))
		if err != nil {
			return
		}
		reply.Type = events.PostRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
