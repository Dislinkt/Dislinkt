package api

import (
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
		user := mapCommandUser(command)
		if err := validator.New().Struct(user); err != nil {
			//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
			//	handler.logger.WarnLogger.Warn(err.Error())
			reply.Type = events.UserServiceNotUpdated
			return
		}
		err := handler.userService.Insert(command.Context, user)
		if err != nil {
			//	handler.logger.WarnLogger.Warn(err.Error())
			reply.Type = events.UserServiceNotUpdated
			return
		}

		//	handler.logger.InfoLogger.Infof("New user registered: {%s}", command.User.Id)
		reply.Type = events.UserServiceUpdated
	case events.RollbackUser:
		err := handler.userService.Delete(mapCommandUser(command))
		if err != nil {
			//	handler.logger.WarnLogger.Warn(err.Error())
			return
		}
		//	handler.logger.InfoLogger.Infof("User deleted: {%s}", command.User.Id)
		reply.Type = events.UserServiceRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
