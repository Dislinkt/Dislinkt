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

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}
