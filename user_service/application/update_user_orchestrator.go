package application

import (
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/user_service/domain"
)

type UpdateUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewUpdateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateUserOrchestrator,
	error) {
	o := &UpdateUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *UpdateUserOrchestrator) Start(user *domain.User) error {
	event := &events.UpdateUserCommand{
		Type: events.UpdateInUser,
		User: events.User{
			Id:       user.Id.String(),
			Name:     user.Name,
			Surname:  user.Surname,
			Username: *user.Username,
		},
	}
	return o.commandPublisher.Publish(event)
}

func (o *UpdateUserOrchestrator) handle(reply *events.UpdateUserReply) {
	command := events.UpdateUserCommand{User: reply.User}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownUpdateCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *UpdateUserOrchestrator) nextCommandType(reply events.UpdateUserReplyType) events.UpdateUserCommandType {
	switch reply {
	case events.UserUpdatedInUser:
		return events.UpdateInPost
	case events.UserNotUpdatedInUser:
		return events.UserUpdateCancelled
	case events.UserRolledBackInUser:
		return events.UserUpdateCancelled
	case events.UserUpdatedInPost:
		return events.UserUpdateSucceeded
	case events.UserNotUpdatedInPost:
		return events.RollbackUpdateInUser
	case events.UserRolledBackInPost:
		return events.RollbackUpdateInUser

	default:
		return events.UnknownUpdateCommand
	}
}
