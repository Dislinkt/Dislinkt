package application

import (
	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/register_user"
	"github.com/dislinkt/user_service/domain"
)

type RegisterUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewRegisterUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*RegisterUserOrchestrator, error) {
	o := &RegisterUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *RegisterUserOrchestrator) Start(user *domain.User) error {
	event := &events.RegisterUserCommand{
		Type: events.UpdateUser,
		User: events.User{
			Id:          "",
			Name:        user.Name,
			Surname:     user.Surname,
			Username:    *user.Username,
			Email:       *user.Email,
			Number:      user.Number,
			Gender:      events.Gender(user.Gender),
			DateOfBirth: user.DateOfBirth,
			Password:    user.Password,
			Biography:   user.Biography,
			Private:     user.Private,
		},
	}
	return o.commandPublisher.Publish(event)
}

func (o *RegisterUserOrchestrator) handle(reply *events.RegisterUserReply) {
	command := events.RegisterUserCommand{User: reply.User}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *RegisterUserOrchestrator) nextCommandType(reply events.RegisterUserReplyType) events.RegisterUserCommandType {
	switch reply {
	case events.UserServiceUpdated:
		return events.UpdateAdditional
	case events.UserServiceNotUpdated:
		return events.CancelRegistration
	case events.AdditionalServiceUpdated:
		return events.UpdateConnectionNode
	case events.AdditionalServiceNotUpdated:
		return events.RollbackUser
	case events.AdditionalServiceRolledBack:
		return events.RollbackUser
	case events.ConnectionsUpdated:
		return events.UpdateAuth
	case events.ConnectionsNotUpdated:
		return events.RollbackAdditional
	case events.ConnectionsRolledBack:
		return events.RollbackAdditional
	case events.AuthUpdated:
		return events.ApproveRegistration
	case events.AuthNotUpdated:
		return events.RollbackConnectionNode
	case events.AuthRolledBack:
		return events.RollbackConnectionNode

	default:
		return events.UnknownCommand
	}
}
