package startup

import (
	"fmt"
	"log"
	"net"

	"github.com/dislinkt/additional_user_service/application"
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/additional_user_service/infrastructure/api"
	"github.com/dislinkt/additional_user_service/infrastructure/persistence"
	"github.com/dislinkt/additional_user_service/startup/config"
	additionalUserProto "github.com/dislinkt/common/proto/additional_user_service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type Server struct {
	config *config.Config
	// tracer otgo.Tracer
	// closer io.Closer
}

func NewServer(config *config.Config) *Server {
	// newTracer, closer := tracer.Init(config.JaegerServiceName)
	// otgo.SetGlobalTracer(newTracer)
	return &Server{
		config: config,
		// tracer: newTracer,
		// closer: closer,
	}
}

const (
	QueueGroup = "additional_user_service"
)

func (server *Server) Start() {
	mongoClient := server.initAdditionalUserClient()
	additionalUserStore := server.initAdditionalUserStore(mongoClient)

	additionalUserService := server.initAdditionalUserService(additionalUserStore)

	additionalUserHandler := server.initAdditionalUserHandler(additionalUserService)

	server.startGrpcServer(additionalUserHandler)
}

func (server *Server) initAdditionalUserClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.AdditionalUserDBHost, server.config.AdditionalUserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initAdditionalUserStore(client *mongo.Client) domain.AdditionalUserStore {
	store := persistence.NewAdditionalUserMongoDBStore(client)
	return store
}

func (server *Server) initAdditionalUserService(store domain.AdditionalUserStore) *application.AdditionalUserService {
	return application.NewAdditionalUserService(store)
}

func (server *Server) initAdditionalUserHandler(service *application.AdditionalUserService) *api.AdditionalUserHandler {
	return api.NewProductHandler(service)
}

func (server *Server) startGrpcServer(additionalUserHandler *api.AdditionalUserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	additionalUserProto.RegisterAdditionalUserServiceServer(grpcServer, additionalUserHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
