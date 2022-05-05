package domain

import (
	pb "github.com/dislinkt/common/proto/connection_service"
)

type ConnectionStore interface {
	Register(userNode *UserNode) (*UserNode, error)
	CreateConnection(baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error)
}
