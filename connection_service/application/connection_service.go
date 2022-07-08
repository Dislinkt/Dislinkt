package application

import (
	"errors"
	"fmt"
	uuid "github.com/gofrs/uuid"

	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/domain"
)

type ConnectionService struct {
	store domain.ConnectionStore
}

func NewConnectionService(store domain.ConnectionStore) *ConnectionService {
	return &ConnectionService{
		store: store,
	}
}

func (service *ConnectionService) Register(userID string, status string) (*domain.UserNode, error) {
	fmt.Println("[ConnectionService Register]")
	if !IsValidUUID(userID) {
		return nil, errors.New("Invalid uuid")
	}
	node := domain.UserNode{UserUID: userID, Status: domain.ProfileStatus(status)}
	return service.store.Register(&node)
}

func (service *ConnectionService) CreateConnection(baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error) {
	fmt.Println("[ConnectionService CreateConnection]")
	if !IsValidUUID(baseUserUuid) {
		return nil, errors.New("Invalid uuid")
	}
	return service.store.CreateConnection(baseUserUuid, connectUserUuid)
}

func (service *ConnectionService) AcceptConnection(requestSenderUser string, requestApprovalUser string) (*pb.NewConnectionResponse, error) {
	fmt.Println("[ConnectionService AcceptConnection")
	return service.store.AcceptConnection(requestSenderUser, requestApprovalUser)
}

func (service *ConnectionService) GetAllConnectionForUser(userUid string) ([]*domain.UserNode, error) {
	if !IsValidUUID(userUid) {
		return nil, errors.New("Invalid uuid")
	}
	return service.store.GetAllConnectionForUser(userUid)
}

func (service *ConnectionService) GetAllConnectionRequestsForUser(userUid string) ([]*domain.UserNode, error) {
	if !IsValidUUID(userUid) {
		return nil, errors.New("Invalid uuid")
	}
	return service.store.GetAllConnectionRequestsForUser(userUid)
}

func (service *ConnectionService) UpdateUser(userUUID string,
	private bool) error {
	if !IsValidUUID(userUUID) {
		return errors.New("Invalid uuid")
	}
	fmt.Println("[ConnectionService UpdateUser")
	return service.store.UpdateUser(userUUID, private)
}

func (service *ConnectionService) BlockUser(currentUserUUID string, blockedUserUUid string) (*pb.BlockedUserStatus, error) {
	fmt.Println("[ConnectionService BlockUser")
	return service.store.BlockUser(currentUserUUID, blockedUserUUid)
}

func (servioe *ConnectionService) GetAllBlockedForCurrentUser(currentUserUUID string) ([]*domain.UserNode, error) {
	fmt.Println("[ConnectionService GetAllBlockedForCurrentUser")
	return servioe.store.GetAllBlockedForCurrentUser(currentUserUUID)
}

func (servioe *ConnectionService) GetAllUserBlockingCurrentUser(currentUserUUID string) ([]*domain.UserNode, error) {
	fmt.Println("[ConnectionService GetAllUserBlockingCurrentUser")
	return servioe.store.GetAllUserBlockingCurrentUser(currentUserUUID)
}

func (servioe *ConnectionService) RecommendUsersByConnection(currentUserUUID string) (users []*domain.UserNode, err error) {
	fmt.Println("[ConnectionService RecommendUsersByConnection")
	return servioe.store.RecommendUsersByConnection(currentUserUUID)
}

func (servioe *ConnectionService) UnblockConnection(currentUser string, blockedUser string) (*pb.BlockedUserStatus, error) {
	fmt.Println("[ConnectionService UnblockConnection")
	return servioe.store.UnblockConnection(currentUser, blockedUser)
}

func (servioe *ConnectionService) InsertField(name string) (string, error) {
	fmt.Println("[ConnectionService InsertField")
	return servioe.store.InsertField(name)
}

func (servioe *ConnectionService) InsertSkill(name string) (string, error) {
	fmt.Println("[ConnectionService InsertSkill")
	return servioe.store.InsertSkill(name)
}

func (servioe *ConnectionService) InsertJobOffer(jobOffer domain.JobOffer) (string, error) {
	fmt.Println("[ConnectionService InsertJobOffer")
	return servioe.store.InsertJobOffer(jobOffer)
}

func (servioe *ConnectionService) InsertSkillToUser(name string, uuid string) (string, error) {
	fmt.Println("[ConnectionService InsertSkillToUser")
	return servioe.store.InsertSkillToUser(name, uuid)
}

func (servioe *ConnectionService) InsertFieldToUser(name string, uuid string) (string, error) {
	fmt.Println("[ConnectionService InsertFieldToUser")
	fmt.Println(name, uuid)
	return servioe.store.InsertFieldToUser(name, uuid)
}

func (servioe *ConnectionService) RecommendJobBySkill(userUid string) (jobs []*domain.JobOffer, err error) {
	fmt.Println("[ConnectionService RecommendJobBySkill")
	return servioe.store.RecommendJobBySkill(userUid)
}

func (servioe *ConnectionService) RecommendJobByField(userUid string) (jobs []*domain.JobOffer, err error) {
	fmt.Println("[ConnectionService RecommendJobByField11")
	return servioe.store.RecommendJobByField(userUid)
}

func (servioe *ConnectionService) CheckIfUsersConnected(userUUID1 string, userUUID2 string) (bool, error) {
	fmt.Println("[ConnectionService CHeckIfUsersConnected")
	return servioe.store.CheckIfUsersConnected(userUUID1, userUUID2)
}

func (servioe *ConnectionService) CheckIfUsersBlocked(userUUID1 string, userUUID2 string) (bool, error) {
	fmt.Println("[ConnectionService CheckIfUsersBlocked")
	return servioe.store.CheckIfUsersBlocked(userUUID1, userUUID2)
}

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}
