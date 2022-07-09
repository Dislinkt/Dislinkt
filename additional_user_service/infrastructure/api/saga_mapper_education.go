package api

import (
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/common/saga/events"
)

func mapAdditionalCommandAddEducation(command *events.AddEducationCommand) *domain.Education {

	educationD := &domain.Education{
		Id:           command.Education.Id,
		School:       command.Education.School,
		Degree:       command.Education.Degree,
		FieldOfStudy: command.Education.FieldOfStudy,
		StartDate:    command.Education.StartDate,
		EndDate:      command.Education.EndDate,
	}
	return educationD
}
