package api

import (
	pb "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/user-service/domain"
)

func mapUser(userD *domain.User) *pb.User {
	userPb := &pb.User{
		Id:      userD.Id,
		Name:    userD.Name,
		Surname: userD.Surname,
	}
	return userPb
}
