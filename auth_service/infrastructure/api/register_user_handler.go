package api

import (
	"fmt"
	"github.com/dislinkt/auth_service/application"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/pquerna/otp/totp"
)

type CreateUserCommandHandler struct {
	userService       *application.UserService
	authService       *application.AuthService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewRegisterUserCommandHandler(userService *application.UserService, authService *application.AuthService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		userService:       userService,
		authService:       authService,
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
	case events.UpdateAuth:
		fmt.Println("UpdateAuth")
		user := mapCommandUser(command)
		if user == nil {
			return
		}

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "NenadBecanovic1234567",
			AccountName: user.Username,
		})

		user.TotpToken = key.Secret()
		uuid, err := handler.userService.Insert(user)

		reply.User.Id = uuid.String()
		if err != nil {
			reply.Type = events.AuthNotUpdated
			return
		}

		errSendingMail := handler.authService.SendActivationMail(user.Username)
		if errSendingMail != nil {
			fmt.Println("ERROR SENDING MAIL")
			reply.Type = events.AuthNotUpdated
			return
		}
		reply.Type = events.AuthUpdated
	case events.RollbackAuth:
		fmt.Println("RollbackAuth")
		err := handler.userService.Delete(mapCommandUser(command))
		if err != nil {
			return
		}
		reply.Type = events.AuthRolledBack
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
