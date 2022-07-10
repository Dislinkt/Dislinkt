package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/additional_user_service/application"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type DeleteEducationCommandHandler struct {
	additionalService *application.AdditionalUserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewDeleteEducationCommandHandler(additionalService *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*DeleteEducationCommandHandler, error) {
	o := &DeleteEducationCommandHandler{
		additionalService: additionalService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *DeleteEducationCommandHandler) handle(command *events.DeleteEducationCommand) {
	reply := events.DeleteEducationReply{
		Education: command.Education,
		UserId:    command.UserId,
	}

	switch command.Type {
	case events.DeleteEducationInAdditional:
		fmt.Println("additional handler add skill")
		education, _ := handler.additionalService.FindUserField(context.TODO(), command.Education.Id, command.UserId)
		_, err := handler.additionalService.DeleteUserEducation(context.TODO(), command.UserId, command.Education.Id)
		if err != nil {
			fmt.Println("additional handler error not added education")
			reply.Type = events.AdditionalServiceEducationNotDeleted
			return
		}

		reply.Education.FieldOfStudy = education.FieldOfStudy
		reply.Type = events.AdditionalServiceEducationDeleted
		fmt.Println("additional handler add success education")
		fmt.Println(reply.Education.FieldOfStudy)
		// reply.Type = events.RegistrationApproved
	case events.RollbackDeleteEducationInAdditional:
		fmt.Println("additional handler-rollback education")
		education, _ := handler.additionalService.FindUserField(context.TODO(), command.Education.Id, command.UserId)
		err, _ := handler.additionalService.CreateEducation(context.TODO(), command.UserId, education)
		if err != nil {
			return
		}
		reply.Type = events.AdditionalEducationDeleteRolledBack
	default:
		reply.Type = events.UnknownDeletedEducationReply
	}

	if reply.Type != events.UnknownDeletedEducationReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
