package api

import (
	"context"
	"fmt"
	logger "github.com/dislinkt/common/logging"

	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"post_service/application"
)

type CreateUserCommandHandler struct {
	postService       *application.PostService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
	logger            *logger.Logger
}

func NewRegisterUserCommandHandler(postService *application.PostService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	logger := logger.InitLogger(context.TODO())
	o := &CreateUserCommandHandler{
		postService:       postService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
		logger:            logger,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateUserCommandHandler) handle(command *events.RegisterUserCommand) {
	reply := events.RegisterUserReply{User: command.User}
	handler.logger.InfoLogger.Infof("SS-UU {%s}", reply.User.Username)

	switch command.Type {
	case events.UpdatePost:
		fmt.Println("posthandler-update")

		err := handler.postService.InsertUser(mapPostCommandUser(command))
		if err != nil {
			handler.logger.WarnLogger.Warnf("SF-UU {%s}", reply.User.Username)
			fmt.Println(err)
			reply.Type = events.PostNotUpdated
			return
		}
		reply.Type = events.PostUpdated
	case events.RollbackPost:
		handler.logger.WarnLogger.Warnf("SF-UU {%s}", reply.User.Username)
		fmt.Println("posthandler-rollback")
		err := handler.postService.DeleteUser(mapPostCommandUser(command))
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
