package startup

import (
	"fmt"
	"github.com/dislinkt/common/interceptor"
	notificationProto "github.com/dislinkt/common/proto/notification_service"
	"github.com/dislinkt/notification_service/application"
	"github.com/dislinkt/notification_service/domain"
	"github.com/dislinkt/notification_service/infrastructure/api"
	"github.com/dislinkt/notification_service/infrastructure/persistence"
	"github.com/dislinkt/notification_service/startup/config"
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
	mongoClient := server.initMongoClient()
	notificationStore := server.initNotificationStore(mongoClient)

	notificationService := server.initNotificationService(notificationStore)
	notificationHandler := server.initNotificationHandler(notificationService)
	server.startGrpcServer(notificationHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.NotificationDBHost, server.config.NotificationDBPort)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func (server *Server) initNotificationStore(client *mongo.Client) domain.NotificationStore {
	store := persistence.NewNotificationMongoDBStore(client)
	return store
}

func (server *Server) initNotificationService(store domain.NotificationStore) *application.NotificationService {
	return application.NewNotificationService(store)
}

func (server *Server) initNotificationHandler(service *application.NotificationService) *api.NotificationHandler {
	return api.NewNotificationHandler(service)
}

func (server *Server) startGrpcServer(notificationHandler *api.NotificationHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	notificationProto.RegisterNotificationServiceServer(grpcServer, notificationHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
