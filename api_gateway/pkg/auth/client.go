package auth

import (
	"fmt"

	"github.com/dislinkt/api_gateway/startup/config"
	pb "github.com/dislinkt/common/proto/auth_service"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.AuthServiceClient
}

func InitServiceClient(c *config.Config) pb.AuthServiceClient {
	// using WithInsecure() because no SSL running
	authEndpoint := fmt.Sprintf("%s:%s", c.AuthHost, c.AuthPort)
	cc, err := grpc.Dial(authEndpoint, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewAuthServiceClient(cc)
}
