package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/additional_user_service/application"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type UpdateEducationCommandHandler struct {
	additionalService *application.AdditionalUserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateEducationCommandHandler(additionalService *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*UpdateEducationCommandHandler, error) {
	o := &UpdateEducationCommandHandler{
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

func (handler *UpdateEducationCommandHandler) handle(command *events.UpdateEducationCommand) {
	reply := events.UpdateEducationReply{
		Education: command.Education,
		UserId:    command.UserId,
	}

	switch command.Type {
	case events.UpdateEducationInAdditional:
		fmt.Println("additional handler update Education")
		education, _ := handler.additionalService.FindUserField(context.TODO(), command.Education.Id, command.UserId)
		education_update := mapAdditionalCommandUpdateEducation(command)
		_, err := handler.additionalService.UpdateUserEducation(context.TODO(), command.UserId, command.Education.Id, education_update)
		if err != nil {
			fmt.Println("additional handler error not updated Education")
			reply.Type = events.AdditionalServiceEducationNotUpdated
			return
		}

		if education.FieldOfStudy != education_update.FieldOfStudy {
			reply.OldFieldName = education.FieldOfStudy
			reply.Type = events.AdditionalServiceEducationUpdated
		} else {
			reply.Type = events.GraphDatabaseEducationUpdated
		}

		fmt.Println("additional handler update success Education")
		fmt.Println(reply.Education.FieldOfStudy)
		fmt.Println(reply.OldFieldName)
		// reply.Type = events.RegistrationApproved
	case events.RollbackEducationUpdateInAdditional:
		fmt.Println("additional handler-rollback Education update")
		education, _ := handler.additionalService.FindUserField(context.TODO(), command.Education.Id, command.UserId)
		err, _ := handler.additionalService.CreateEducation(context.TODO(), command.UserId, education)
		if err != nil {
			return
		}
		reply.Type = events.AdditionalEducationUpdateRolledBack
	default:
		reply.Type = events.UnknownUpdatedEducationReply
	}

	if reply.Type != events.UnknownUpdatedEducationReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
