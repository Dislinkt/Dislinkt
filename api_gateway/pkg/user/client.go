package user

import (
	"fmt"

	"github.com/dislinkt/api_gateway/startup/config"
	pb "github.com/dislinkt/common/proto/user_service"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.UserServiceClient
}

func InitServiceClient(c *config.Config) pb.UserServiceClient {
	// using WithInsecure() because no SSL running
	userEndpoint := fmt.Sprintf("%s:%s", c.UserHost, c.UserPort)
	cc, err := grpc.Dial(userEndpoint, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewUserServiceClient(cc)
}
