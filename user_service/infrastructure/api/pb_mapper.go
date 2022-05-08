package api

import (
	"time"

	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/user_service/domain"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapUser(userD *domain.User) *pb.User {
	userPb := &pb.User{
		Id:          userD.Id.String(),
		Name:        userD.Name,
		Surname:     userD.Surname,
		Username:    *userD.Username,
		Email:       *userD.Email,
		Number:      userD.Number,
		Gender:      mapGenderToPb(userD.Gender),
		DateOfBirth: userD.DateOfBirth,
		Password:    userD.Password,
		UserRole:    mapUserRole(userD.UserRole),
		Biography:   userD.Biography,
		Blocked:     userD.Blocked,
		CreatedAt:   timestamppb.New(userD.CreatedAt),
		UpdatedAt:   timestamppb.New(userD.UpdatedAt),
		Private:     userD.Private,
	}
	return userPb
}

func mapNewUser(userPb *pb.NewUser) *domain.User {
	userD := &domain.User{
		Id:          uuid.NewV4(),
		Name:        userPb.Name,
		Surname:     userPb.Surname,
		Username:    &userPb.Username,
		Email:       &userPb.Email,
		Number:      userPb.Number,
		Gender:      mapGenderToDomain(userPb.Gender),
		DateOfBirth: userPb.DateOfBirth,
		Password:    userPb.Password,
		UserRole:    domain.Regular,
		Biography:   userPb.Biography,
		Blocked:     false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Private:     userPb.Private,
	}
	return userD
}

func mapGenderToPb(gender domain.Gender) pb.Gender {
	switch gender {
	case domain.Empty:
		return pb.Gender_Empty
	case domain.Male:
		return pb.Gender_Male
	case domain.Female:
		return pb.Gender_Female
	}
	return pb.Gender_Empty
}

func mapGenderToDomain(gender pb.Gender) domain.Gender {
	switch gender {
	case pb.Gender_Empty:
		return domain.Empty
	case pb.Gender_Male:
		return domain.Male
	case pb.Gender_Female:
		return domain.Female
	}
	return domain.Empty
}

func mapUserRole(role domain.Role) pb.UserRole {
	switch role {
	case domain.Agent:
		return pb.UserRole_Agent
	case domain.Admin:
		return pb.UserRole_Admin
	}
	return pb.UserRole_Regular
}
