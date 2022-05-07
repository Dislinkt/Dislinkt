package domain

import (
	pb "github.com/dislinkt/common/proto/connection_service"
)

type ConnectionStore interface {
	Register(userNode *UserNode) (*UserNode, error)
	CreateConnection(baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error)
	AcceptConnection(requestSenderUser string, requestApprovalUser string) (*pb.NewConnectionResponse, error)
	GetAllConnectionForUser(userUid string) ([]*UserNode, error)
}
