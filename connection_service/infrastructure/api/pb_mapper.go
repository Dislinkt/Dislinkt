package api

import (
	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/domain"
	"google.golang.org/genproto/googleapis/type/date"
)

func mapUserConn(userConn *domain.UserNode) *pb.User {
	userConnPb := &pb.User{
		UserID: userConn.UserUID,
		Status: string(userConn.Status),
	}

	return userConnPb
}

func mapJobOffer(jobOffer *domain.JobOffer) *pb.JobOffer {

	jobOfferPb := &pb.JobOffer{
		Id:            jobOffer.Id,
		Position:      jobOffer.Position,
		Preconditions: jobOffer.Preconditions,
		Duration:      jobOffer.Duration,
		Location:      jobOffer.Location,
		Title:         jobOffer.Title,
		Field:         jobOffer.Field,
	}
	return jobOfferPb
}

func mapJobOfferPb(jobOfferPb *pb.JobOffer) *domain.JobOffer {

	jobOffer := &domain.JobOffer{
		Id:            jobOfferPb.Id,
		Position:      jobOfferPb.Position,
		Preconditions: jobOfferPb.Preconditions,
		DatePosted:    date.Date{},
		Duration:      jobOfferPb.Duration,
		Location:      jobOfferPb.Location,
		Title:         jobOfferPb.Title,
		Field:         jobOfferPb.Field,
	}
	return jobOffer
}
