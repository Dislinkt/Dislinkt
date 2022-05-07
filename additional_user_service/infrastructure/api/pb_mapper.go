package api

import (
	"time"

	"github.com/dislinkt/additional_user_service/domain"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EDUCATION

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

func mapEducations(educations *map[string]domain.Education) []*pb.Education {
	if educations == nil {
		return nil
	}
	var educationsPb []*pb.Education
	for _, education := range *educations {
		educationsPb = append(educationsPb, mapEducation(&education))
	}
	return educationsPb
}

// POSITION

func mapNewPosition(positionPb *pb.NewPosition) *domain.Position {
	startDate, _ := time.Parse("2006-01-02", positionPb.StartDate)
	endDate, _ := time.Parse("2006-01-02", positionPb.EndDate)

	positionD := &domain.Position{
		Title:       positionPb.Title,
		CompanyName: positionPb.CompanyName,
		Industry:    positionPb.Industry,
		StartDate:   primitive.NewDateTimeFromTime(startDate),
		EndDate:     primitive.NewDateTimeFromTime(endDate),
		Current:     positionPb.Current,
	}
	return positionD
}

func mapPosition(position *domain.Position) *pb.Position {
	positionPb := &pb.Position{
		Id:          position.Id.Hex(),
		Title:       position.Title,
		CompanyName: position.CompanyName,
		Industry:    position.Industry,
		StartDate:   timestamppb.New(position.StartDate.Time()),
		EndDate:     timestamppb.New(position.EndDate.Time()),
		Current:     position.Current,
	}
	return positionPb
}

func mapPositions(positions *map[string]domain.Position) []*pb.Position {
	if positions == nil {
		return nil
	}
	var positionsPb []*pb.Position
	for _, position := range *positions {
		positionsPb = append(positionsPb, mapPosition(&position))
	}
	return positionsPb
}

// SKILL

func mapNewSkill(skillPb *pb.NewSkill) *domain.Skill {
	skillD := &domain.Skill{
		Name: skillPb.Name,
	}
	return skillD
}

func mapSkill(skill *domain.Skill) *pb.Skill {
	skillPb := &pb.Skill{
		Id:   skill.Id.Hex(),
		Name: skill.Name,
	}
	return skillPb
}

func mapSkills(skills *map[string]domain.Skill) []*pb.Skill {
	if skills == nil {
		return nil
	}
	var skillsPb []*pb.Skill
	for _, skill := range *skills {
		skillsPb = append(skillsPb, mapSkill(&skill))
	}
	return skillsPb
}

// INTEREST

func mapNewInterest(interestPb *pb.NewInterest) *domain.Interest {
	interestD := &domain.Interest{
		Name: interestPb.Name,
	}
	return interestD
}

func mapInterest(skill *domain.Interest) *pb.Interest {
	skillPb := &pb.Interest{
		Id:   skill.Id.Hex(),
		Name: skill.Name,
	}
	return skillPb
}

func mapInterests(interests *map[string]domain.Interest) []*pb.Interest {
	if interests == nil {
		return nil
	}
	var interestsPb []*pb.Interest
	for _, interest := range *interests {
		interestsPb = append(interestsPb, mapInterest(&interest))
	}
	return interestsPb
}
