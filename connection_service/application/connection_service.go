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
	node := domain.UserNode{userID, domain.ProfileStatus(status)}
	return service.store.Register(&node)
}

func (service *ConnectionService) CreateConnection(baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error) {
	fmt.Println("[ConnectionService CreateConnection]")
	return service.store.CreateConnection(baseUserUuid, connectUserUuid)
}
