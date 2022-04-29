package api

import (
	"time"

	"github.com/dislinkt/additional_user_service/domain"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapNewEducation(educationPb *pb.NewEducation) *domain.Education {
	startDate, _ := time.Parse("2006-01-02", educationPb.StartDate)
	endDate, _ := time.Parse("2006-01-02", educationPb.EndDate)

	educationD := &domain.Education{
		Degree:       educationPb.Degree,
		School:       educationPb.School,
		FieldOfStudy: educationPb.FieldOfStudy,
		StartDate:    primitive.NewDateTimeFromTime(startDate),
		EndDate:      primitive.NewDateTimeFromTime(endDate),
	}
	return educationD
}

func mapEducation(education *domain.Education) *pb.Education {
	educationPb := &pb.Education{
		Id:           education.Id.Hex(),
		School:       education.School,
		Degree:       education.Degree,
		FieldOfStudy: education.FieldOfStudy,
		StartDate:    timestamppb.New(education.StartDate.Time()),
		EndDate:      timestamppb.New(education.EndDate.Time()),
	}
	return educationPb
}

func mapEducations(educations *[]domain.Education) []*pb.Education {
	var educationsPb []*pb.Education
	for _, education := range *educations {
		educationsPb = append(educationsPb, mapEducation(&education))
	}
	return educationsPb
}
