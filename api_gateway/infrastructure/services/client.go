package services

import (
	connectionGw "github.com/dislinkt/common/proto/connection_service"
	postGw "github.com/dislinkt/common/proto/post_service"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//func NewUserClient(address string) userGw.UserServiceClient {
//	conn, err := getConnection(address)
//	if err != nil {
//		log.Fatalf("Failed to start gRPC connection to User service: %v", err)
//	}
//	return userGw.NewUserServiceClient(conn)
//}

func NewPostClient(address string) postGw.PostServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to User service: %v", err)
	}
	return postGw.NewPostServiceClient(conn)
}

func NewConnectionClient(address string) connectionGw.ConnectionServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to User service: %v", err)
	}
	return connectionGw.NewConnectionServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
