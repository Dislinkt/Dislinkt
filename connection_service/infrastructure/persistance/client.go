package persistance

import (
	"fmt"
	eventGw "github.com/dislinkt/common/proto/event_service"
	notificationGw "github.com/dislinkt/common/proto/notification_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GetClient(uri, username, password string) (*neo4j.Driver, error) {

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &driver, nil //TODO: ref driver ?
}

func UserClient(address string) userGw.UserServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to User service: %v", err)
	}
	return userGw.NewUserServiceClient(conn)
}

func NotificationClient(address string) notificationGw.NotificationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Notification service: %v", err)
	}
	return notificationGw.NewNotificationServiceClient(conn)
}

func EventClient(address string) eventGw.EventServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Event service: %v", err)
	}
	return eventGw.NewEventServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
