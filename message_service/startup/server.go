package startup

import (
	"fmt"
	"github.com/dislinkt/common/interceptor"
	messageProto "github.com/dislinkt/common/proto/message_service"
	"github.com/dislinkt/common/tracer"
	"github.com/dislinkt/message_service/application"
	"github.com/dislinkt/message_service/domain"
	"github.com/dislinkt/message_service/infrastructure/api"
	"github.com/dislinkt/message_service/infrastructure/persistence"
	"github.com/dislinkt/message_service/startup/config"
	otgo "github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}

func (server *Server) Start() {
	tracer, _ := tracer.Init("message_service")
	otgo.SetGlobalTracer(tracer)
	mongoClient := server.initMongoClient()
	messageStore := server.initMessageStore(mongoClient)

	messageService := server.initMessageService(messageStore)
	messageHandler := server.initMessageHandler(messageService)
	server.startGrpcServer(messageHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.MessageDBHost, server.config.MessageDBPort)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func (server *Server) initMessageStore(client *mongo.Client) domain.MessageStore {
	store := persistence.NewMessagesMongoDBStore(client)
	return store
}

func (server *Server) initMessageService(store domain.MessageStore) *application.MessageService {
	return application.NewMessageService(store)
}

func (server *Server) initMessageHandler(service *application.MessageService) *api.MessageHandler {
	return api.NewMessageHandler(service)
}

func (server *Server) startGrpcServer(messageHandler *api.MessageHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	messageProto.RegisterMessageServiceServer(grpcServer, messageHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
