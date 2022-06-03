package application

import (
	"fmt"

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
	node := domain.UserNode{UserUID: userID, Status: domain.ProfileStatus(status)}
	return service.store.Register(&node)
}

func (service *ConnectionService) CreateConnection(baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error) {
	fmt.Println("[ConnectionService CreateConnection]")
	return service.store.CreateConnection(baseUserUuid, connectUserUuid)
}

func (service *ConnectionService) AcceptConnection(requestSenderUser string, requestApprovalUser string) (*pb.NewConnectionResponse, error) {
	fmt.Println("[ConnectionService AcceptConnection")
	return service.store.AcceptConnection(requestSenderUser, requestApprovalUser)
}

func (service *ConnectionService) GetAllConnectionForUser(userUid string) ([]*domain.UserNode, error) {
	return service.store.GetAllConnectionForUser(userUid)
}

func (service *ConnectionService) GetAllConnectionRequestsForUser(userUid string) ([]*domain.UserNode, error) {
	return service.store.GetAllConnectionRequestsForUser(userUid)
}

func (service *ConnectionService) UpdateUser(userUUID string,
	private bool) error {
	fmt.Println("[ConnectionService UpdateUser")
	return service.store.UpdateUser(userUUID, private)
}
