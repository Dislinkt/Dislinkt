package startup

import (
	"fmt"
	"github.com/dislinkt/common/interceptor"
	postProto "github.com/dislinkt/common/proto/post_service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	"post_service/application"
	"post_service/domain"
	"post_service/infrastructure/api"
	"post_service/infrastructure/persistence"
	"post_service/startup/config"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	postStore := server.initPostStore(mongoClient)

	postService := server.initPostService(postStore)
	postHandler := server.initPostHandler(postService)
	server.startGrpcServer(postHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func (server *Server) initPostStore(client *mongo.Client) domain.PostStore {
	store := persistence.NewPostMongoDBStore(client)
	return store
}

func (server *Server) initPostService(store domain.PostStore) *application.PostService {
	return application.NewPostService(store)
}

func (server *Server) initPostHandler(service *application.PostService) *api.PostHandler {
	return api.NewPostHandler(service)
}

func (server *Server) startGrpcServer(postHandler *api.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessibleRoles(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	postProto.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
