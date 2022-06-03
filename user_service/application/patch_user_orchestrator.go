package application

import (
	saga "github.com/dislinkt/common/saga/messaging"
	events "github.com/dislinkt/common/saga/patch_user"
	"github.com/dislinkt/user_service/domain"
)

type PatchUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewPatchUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*PatchUserOrchestrator, error) {
	o := &PatchUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *PatchUserOrchestrator) Start(user *domain.User) error {
	event := &events.PatchUserCommand{
		Type: events.PatchUserInUser,
		User: events.User{
			Id:       "",
			Private:  user.Private,
			Username: *user.Username,
		},
	}
	return o.commandPublisher.Publish(event)
}

func (o *PatchUserOrchestrator) handle(reply *events.PatchUserReply) {
	command := events.PatchUserCommand{User: reply.User}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownPatchCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *PatchUserOrchestrator) nextCommandType(reply events.PatchUserReplyType) events.PatchUserCommandType {
	switch reply {
	case events.PatchedUserInUser:
		return events.PatchUserInConnection
	case events.PatchFailedInUser:
		return events.CancelPatch
	case events.PatchFailedInConnection:
		return events.RollbackPatchInUser
	case events.PatchedInConnection:
		return events.PatchSucceeded

	default:
		return events.UnknownPatchCommand
	}
}
