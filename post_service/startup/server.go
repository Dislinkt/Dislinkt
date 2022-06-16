package startup

import (
	"context"
	"fmt"
	logger "github.com/dislinkt/common/logging"
	"log"
	"net"

	"github.com/dislinkt/common/interceptor"
	postProto "github.com/dislinkt/common/proto/post_service"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/common/saga/messaging/nats"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"post_service/application"
	"post_service/domain"
	"post_service/infrastructure/api"
	"post_service/infrastructure/persistence"
	"post_service/startup/config"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
}

func NewServer(config *config.Config) *Server {
	logger := logger.InitLogger(context.TODO())
	return &Server{
		config: config,
		logger: logger,
	}
}

const (
	QueueGroupRegister = "post_service_register"
	QueueGroupUpdate   = "post_service_update"
)

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	postStore := server.initPostStore(mongoClient)

	postService := server.initPostService(postStore)
	postHandler := server.initPostHandler(postService)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroupRegister)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(postService, replyPublisher, commandSubscriber)

	updateCommandSubscriber := server.initSubscriber(server.config.UpdateUserCommandSubject, QueueGroupUpdate)
	updateReplyPublisher := server.initPublisher(server.config.UpdateUserReplySubject)
	server.initUpdateUserHandler(postService, updateReplyPublisher, updateCommandSubscriber)

	server.startGrpcServer(postHandler)
	server.logger.InfoLogger.Infof("SS")
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		server.logger.ErrorLogger.Error("IC")
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

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initRegisterUserHandler(service *application.PostService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewRegisterUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initUpdateUserHandler(service *application.PostService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewUpdateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) startGrpcServer(postHandler *api.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.logger.ErrorLogger.Error("FTL")
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	postProto.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.logger.ErrorLogger.Error("FTS")
		log.Fatalf("failed to serve: %s", err)
	}
}
