package domain

import (
	pb "github.com/dislinkt/common/proto/connection_service"
)

type ConnectionStore interface {
	Register(userNode *UserNode) (*UserNode, error)
	CreateConnection(baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error)
	AcceptConnection(requestSenderUser string, requestApprovalUser string) (*pb.NewConnectionResponse, error)
	GetAllConnectionForUser(userUid string) ([]*UserNode, error)
	GetAllConnectionRequestsForUser(userUid string) ([]*UserNode, error)
	UpdateUser(userUUID string, private bool) error
	BlockUser(currentUserUUID string, blockedUserUUID string) (*pb.BlockedUserStatus, error)
	GetAllBlockedForCurrentUser(currentUserUUID string) ([]*UserNode, error)
	GetAllUserBlockingCurrentUser(currentUserUUID string) ([]*UserNode, error)
	RecommendUsersByConnection(currentUserUUID string) (users []*UserNode, err error)
}
