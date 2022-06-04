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
		Type: events.UpdateInPost,
		User: events.User{
			Id:          user.Id.String(),
			Name:        user.Name,
			Surname:     user.Surname,
			Username:    *user.Username,
			Number:      user.Number,
			DateOfBirth: user.DateOfBirth,
			Gender:      events.Gender(user.Gender),
			Biography:   user.Biography,
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
	case events.UserUpdatedInPost:
		return events.UpdateInUser
	case events.UserNotUpdatedInPost:
		return events.UserUpdateCancelled
	case events.UserRolledBackInPost:
		return events.UserUpdateCancelled
	case events.UserUpdatedInUser:
		return events.UserUpdateSucceeded
	case events.UserNotUpdatedInUser:
		return events.RollbackUpdateInPost
	case events.UserRolledBackInUser:
		return events.RollbackUpdateInPost

	default:
		return events.UnknownUpdateCommand
	}
}
