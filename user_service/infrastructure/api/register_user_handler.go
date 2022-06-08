package api

import (
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/user_service/application"
	"github.com/go-playground/validator/v10"
)

type CreateUserCommandHandler struct {
	userService       *application.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(userService *application.UserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		userService:       userService,
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
	case events.UpdateUser:
		fmt.Println("userhandler-update")
		fmt.Println(command.User)
		user := mapCommandUser(command)
		if err := validator.New().Struct(user); err != nil {
			//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
			reply.Type = events.UserServiceNotUpdated
			return
		}
		err := handler.userService.Insert(user)
		if err != nil {
			reply.Type = events.UserServiceNotUpdated
			return
		}
		reply.Type = events.UserServiceUpdated
	case events.RollbackUser:
		fmt.Println("userhandler-rollback")
		err := handler.userService.Delete(mapCommandUser(command))
		if err != nil {
			return
		}
		reply.Type = events.UserServiceRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
