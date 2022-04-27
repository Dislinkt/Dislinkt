package services

import (
	"log"

	userGw "github.com/dislinkt/common/proto/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// trebalo je za neku sagu
func NewUserClient(address string) userGw.UserServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to User service: %v", err)
	}
	return userGw.NewUserServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
