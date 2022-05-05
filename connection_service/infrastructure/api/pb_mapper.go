package api

import (
	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/domain"
)

func mapUserConn(userConn *domain.UserNode) *pb.User {
	userConnPb := &pb.User{
		UserID: userConn.UserUID,
		Status: string(userConn.Status),
	}

	return userConnPb
}
