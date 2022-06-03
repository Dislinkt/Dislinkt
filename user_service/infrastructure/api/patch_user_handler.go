package api

import (
	"fmt"

	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/patch_user"
	"github.com/dislinkt/user_service/application"
)

type PatchUserCommandHandler struct {
	userService       *application.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewPatchUserCommandHandler(userService *application.UserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*PatchUserCommandHandler, error) {
	o := &PatchUserCommandHandler{
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

func (handler *PatchUserCommandHandler) handle(command *events.PatchUserCommand) {
	reply := events.PatchUserReply{User: command.User}

	switch command.Type {
	case events.PatchUserInUser:
		fmt.Println("user handler-patch")
		fmt.Println(command.User)
		var paths []string
		paths = append(paths, "private")
		dbUser, err := handler.userService.PatchUser(paths, mapPatchUser(command.User), command.User.Username)
		if err != nil {
			fmt.Println(err)
			reply.Type = events.PatchFailedInUser
			return
		}
		reply.User.Id = dbUser.Id.String()
		reply.Type = events.PatchedUserInUser

	default:
		reply.Type = events.UnknownPatchReply
	}

	if reply.Type != events.UnknownPatchReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
