package api

import (
	pb "github.com/dislinkt/common/proto/connection_service"
	pbEvent "github.com/dislinkt/common/proto/event_service"
	"github.com/dislinkt/connection_service/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

func mapUserConn(userConn *domain.UserNode) *pb.User {
	userConnPb := &pb.User{
		UserID: userConn.UserUID,
		Status: string(userConn.Status),
	}

	return userConnPb
}

func mapJobOffer(jobOffer *domain.JobOffer) *pb.JobOffer {

	duration := strconv.Itoa(jobOffer.Duration)

	jobOfferPb := &pb.JobOffer{
		Id:            jobOffer.Id,
		Position:      jobOffer.Position,
		Preconditions: jobOffer.Preconditions,
		Duration:      duration,
		Location:      jobOffer.Location,
		Title:         jobOffer.Title,
		Field:         jobOffer.Field,
		Description:   jobOffer.Description,
		DatePosted:    timestamppb.New(jobOffer.DatePosted),
	}
	return jobOfferPb
}

func mapJobOfferPb(jobOfferPb *pb.JobOffer) *domain.JobOffer {
	dur, _ := strconv.Atoi(jobOfferPb.Duration)

	jobOffer := &domain.JobOffer{
		Id:            jobOfferPb.Id,
		Position:      jobOfferPb.Position,
		Preconditions: jobOfferPb.Preconditions,
		DatePosted:    jobOfferPb.DatePosted.AsTime(),
		Duration:      dur,
		Location:      jobOfferPb.Location,
		Title:         jobOfferPb.Title,
		Field:         jobOfferPb.Field,
		Description:   jobOfferPb.Description,
	}

	return jobOffer
}

func mapEventForUserPrivacyChange(userId string, isPrivate bool) *pbEvent.NewEvent {
	eventPb := &pbEvent.NewEvent{
		UserId: userId,
	}
	if isPrivate {
		eventPb.Description = "Set account to private."
	} else {
		eventPb.Description = "Set account to public."
	}
	return eventPb
}
