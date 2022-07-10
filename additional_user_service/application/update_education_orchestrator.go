package application

import (
	"fmt"
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
)

type UpdateEducationOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewUpdateEducationOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateEducationOrchestrator,
	error) {
	o := &UpdateEducationOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *UpdateEducationOrchestrator) Start(userId string, educationId string, education *domain.Education) error {
	event := &events.UpdateEducationCommand{
		Type: events.UpdateEducationInAdditional,
		Education: events.EducationUpdate{
			Id:           educationId,
			Degree:       education.Degree,
			FieldOfStudy: education.FieldOfStudy,
			StartDate:    education.StartDate,
			EndDate:      education.EndDate,
			School:       education.School,
		},
		UserId:       userId,
		OldFieldName: "",
	}
	return o.commandPublisher.Publish(event)
}

func (o *UpdateEducationOrchestrator) handle(reply *events.UpdateEducationReply) {
	command := events.UpdateEducationCommand{Education: reply.Education, UserId: reply.UserId, OldFieldName: reply.OldFieldName}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownUpdateEducationCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *UpdateEducationOrchestrator) nextCommandType(reply events.UpdateEducationReplyType) events.UpdateEducationCommandType {
	switch reply {
	case events.AdditionalServiceEducationUpdated:
		fmt.Println("proslo u additional skill update")
		return events.UpdateEducationInGraph
	case events.GraphDatabaseEducationNotUpdated:
		return events.RollbackEducationUpdateInAdditional
	case events.AdditionalEducationUpdateRolledBack:
		return events.CancelEducationUpdate
	case events.AdditionalServiceEducationNotUpdated:
		return events.CancelEducationUpdate
	case events.GraphDatabaseEducationUpdated:
		fmt.Println("proslo u graf skill update")
		return events.ApproveEducationUpdating
	default:
		return events.UnknownUpdateEducationCommand
	}
}
