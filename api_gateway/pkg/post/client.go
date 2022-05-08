package post

import (
	"fmt"
	"github.com/dislinkt/api_gateway/startup/config"
	pb "github.com/dislinkt/common/proto/post_service"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.PostServiceClient
}

func InitServiceClient(c *config.Config) pb.PostServiceClient {
	// using WithInsecure() because no SSL running
	postEndpoint := fmt.Sprintf("%s:%s", c.PostHost, c.PostPort)
	cc, err := grpc.Dial(postEndpoint, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewPostServiceClient(cc)
}
