package api

import (
	"github.com/dislinkt/common/saga/events"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/domain"
)

func mapPostCommandCreateJob(command *events.CreateJobOfferCommand) *domain.JobOffer {

	jobOfferD := &domain.JobOffer{
		Id:            primitive.NewObjectID(),
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
