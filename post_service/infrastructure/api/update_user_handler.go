package api

import (
	"context"
	"fmt"

	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"post_service/application"
)

type UpdateUserCommandHandler struct {
	postService       *application.PostService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateUserCommandHandler(postService *application.PostService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*UpdateUserCommandHandler, error) {
	o := &UpdateUserCommandHandler{
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

func (handler *UpdateUserCommandHandler) handle(command *events.UpdateUserCommand) {
	reply := events.UpdateUserReply{User: command.User}

	switch command.Type {
	case events.UpdateInPost:
		fmt.Println("post handler-update")

		err := handler.postService.UpdateUser(context.TODO(), mapPostCommandUpdateUser(command))
		if err != nil {
			//reply.User =
			reply.Type = events.UserNotUpdatedInPost
			return
		}
		reply.Type = events.UserUpdatedInPost
	default:
		reply.Type = events.UnknownUpdateReply
	}

	if reply.Type != events.UnknownUpdateReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
