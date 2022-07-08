package api

import (
	"github.com/dislinkt/common/saga/events"
	"github.com/dislinkt/connection_service/domain"
)

func mapConnectionCommandCreateJob(command *events.CreateJobOfferCommand) *domain.JobOffer {

	jobOfferD := &domain.JobOffer{
		Id:            command.JobOffer.Id.Hex(),
		Position:      command.JobOffer.Position,
		Description:   command.JobOffer.Description,
		Preconditions: command.JobOffer.Preconditions,
		DatePosted:    command.JobOffer.DatePosted,
		Duration:      command.JobOffer.Duration,
		Location:      command.JobOffer.Location,
		Title:         command.JobOffer.Title,
		Field:         command.JobOffer.Field,
	}
	return jobOfferD
}
