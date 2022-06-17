package api

import (
	"context"
	logger "github.com/dislinkt/common/logging"
	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/patch_user"
	"github.com/dislinkt/user_service/application"
	"github.com/go-playground/validator/v10"
)

type PatchUserCommandHandler struct {
	userService       *application.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
	logger            *logger.Logger
}

func NewPatchUserCommandHandler(userService *application.UserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*PatchUserCommandHandler, error) {
	logger := logger.InitLogger(context.TODO())
	o := &PatchUserCommandHandler{
		userService:       userService,
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

func (handler *PatchUserCommandHandler) handle(command *events.PatchUserCommand) {
	reply := events.PatchUserReply{User: command.User}
	handler.logger.InfoLogger.Infof("SS-PU {%s}", reply.User.Username)

	switch command.Type {
	case events.PatchUserInUser:
		var paths []string
		paths = append(paths, "private")
		user := mapPatchUser(command.User)
		if err := validator.New().Struct(user); err != nil {
			handler.logger.WarnLogger.Warnf("SF-PU {%s}", reply.User.Username)
			//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
			//	handler.logger.WarnLogger.Warn(err.Error())
			reply.Type = events.PatchFailedInUser
			return
		}
		dbUser, err := handler.userService.PatchUser(paths, user, command.User.Username)
		if err != nil {
			handler.logger.WarnLogger.Warnf("SF-PU {%s}", reply.User.Username)
			//	handler.logger.WarnLogger.Warn(err.Error())
			reply.Type = events.PatchFailedInUser
			return
		}
		reply.User.Id = dbUser.Id.String()
		//	handler.logger.InfoLogger.Infof("User updated: {%s}", dbUser.Id.String())
		reply.Type = events.PatchedUserInUser

	default:
		reply.Type = events.UnknownPatchReply
	}

	if reply.Type != events.UnknownPatchReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
