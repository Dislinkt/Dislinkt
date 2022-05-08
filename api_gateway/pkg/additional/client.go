package additional

import (
	"fmt"
	"github.com/dislinkt/api_gateway/startup/config"
	pb "github.com/dislinkt/common/proto/additional_user_service"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.AdditionalUserServiceClient
}

func InitServiceClient(c *config.Config) pb.AdditionalUserServiceClient {
	// using WithInsecure() because no SSL running
	userEndpoint := fmt.Sprintf("%s:%s", c.UserHost, c.UserPort)
	cc, err := grpc.Dial(userEndpoint, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewAdditionalUserServiceClient(cc)
}
