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
	UnblockConnection(currentUser string, blockedUser string) (*pb.BlockedUserStatus, error)
	InsertField(name string) (string, error)
	InsertSkill(name string) (string, error)
	InsertJobOffer(jobOffer JobOffer) (string, error)
	InsertSkillToUser(name string, uuid string) (string, error)
	InsertFieldToUser(name string, uuid string) (string, error)
	RecommendJobBySkill(userUid string) (jobs []*JobOffer, err error)
	RecommendJobByField(userUid string) (jobs []*JobOffer, err error)
	DeleteAllSkills() (res string, err error)
	DeleteAllFields() (res string, err error)
	DeleteSkillForUser(userId string, skillName string) (res string, err error)
	CheckIfUsersConnected(uuid1 string, uuid2 string) (isConnected bool, err error)
	CheckIfUsersBlocked(uuid1 string, uuid2 string) (isBlocked bool, err error)
	UpdateSkillForUser(userUUID string, skillNameOld string, skillNameNew string) (res string, err error)
	UpdateFieldForUser(userUUID string, fieldNameOld string, fieldNameNew string) (res string, err error)
	CheckIfDatabaseFilled() (isFilled bool, err error)
}
