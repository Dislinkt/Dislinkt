package api

import (
	"fmt"
	"time"

	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/user_service/domain"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapUser(userD *domain.User) *pb.User {
	id := userD.Id.String()

	links := &pb.Links{
		User:        "/user/" + id,
		Interests:   "/user/" + id + "/interest",
		Skills:      "/user/" + id + "/skill",
		Positions:   "/user/" + id + "/position",
		Educations:  "/user/" + id + "/education",
		Posts:       "/post/" + id,
		Connections: "/connection/user/" + id,
		Feed:        "/user/" + id + "/feed",
	}

	userPb := &pb.User{
		Id:          userD.Id.String(),
		Name:        userD.Name,
		Surname:     userD.Surname,
		Username:    *userD.Username,
		Email:       *userD.Email,
		Number:      userD.Number,
		Gender:      mapGenderToPb(userD.Gender),
		DateOfBirth: userD.DateOfBirth,
		UserRole:    mapUserRole(userD.UserRole),
		Biography:   userD.Biography,
		Blocked:     userD.Blocked,
		CreatedAt:   timestamppb.New(userD.CreatedAt),
		UpdatedAt:   timestamppb.New(userD.UpdatedAt),
		Private:     userD.Private,
		Links:       links,
	}
	return userPb
}

func mapNewUser(userPb *pb.NewUser) *domain.User {
	fmt.Println("PROTO" + userPb.Password)
	userD := &domain.User{
		Id:        uuid.NewV4(),
		Name:      userPb.Name,
		Surname:   userPb.Surname,
		Username:  &userPb.Username,
		Email:     &userPb.Email,
		UserRole:  domain.Regular,
		Blocked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Private:   userPb.Private,
		Password:  userPb.Password,
	}
	fmt.Println("DOMAIN: " + userD.Password)
	return userD
}

func mapUpdateUser(userPb *pb.UpdateUser) *domain.User {
	userD := &domain.User{
		Id:          uuid.NewV4(),
		Name:        userPb.Name,
		Surname:     userPb.Surname,
		Username:    &userPb.Username,
		Number:      userPb.Number,
		Gender:      mapGenderToDomain(userPb.Gender),
		DateOfBirth: userPb.DateOfBirth,
		UserRole:    domain.Regular,
		Biography:   userPb.Biography,
		UpdatedAt:   time.Now(),
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
