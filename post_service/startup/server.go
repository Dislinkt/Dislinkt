package startup

import (
	"fmt"
	"github.com/dislinkt/common/tracer"
	otgo "github.com/opentracing/opentracing-go"
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
}

func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}

const (
	QueueGroupRegister  = "post_service_register"
	QueueGroupUpdate    = "post_service_update"
	QueueGroupCreateJob = "post_service_create_job"
)

func (server *Server) Start() {
	tracer, _ := tracer.Init("post_service")
	otgo.SetGlobalTracer(tracer)

	mongoClient := server.initMongoClient()
	postStore := server.initPostStore(mongoClient)

	fmt.Println(server.config.CreateJobOfferCommandSubject + "start post_service")
	commandPublisher := server.initPublisher(server.config.CreateJobOfferCommandSubject)
	replySubscriber := server.initSubscriber(server.config.CreateJobOfferReplySubject, QueueGroupCreateJob)
	createJobOfferOrchestrator := server.initCreateJobOfferOrchestrator(commandPublisher, replySubscriber)

	postService := server.initPostService(postStore, createJobOfferOrchestrator)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroupRegister)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(postService, replyPublisher, commandSubscriber)

	updateCommandSubscriber := server.initSubscriber(server.config.UpdateUserCommandSubject, QueueGroupUpdate)
	updateReplyPublisher := server.initPublisher(server.config.UpdateUserReplySubject)
	server.initUpdateUserHandler(postService, updateReplyPublisher, updateCommandSubscriber)

	createJobOfferSubscriber := server.initSubscriber(server.config.CreateJobOfferCommandSubject, QueueGroupCreateJob)
	createJobOfferPublisher := server.initPublisher(server.config.CreateJobOfferReplySubject)
	server.initCreateJobOfferHandler(postService, createJobOfferPublisher, createJobOfferSubscriber)

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

func (server *Server) initPostService(store domain.PostStore, createJobOfferOrchestrator *application.CreateJobOfferOrchestrator) *application.PostService {
	return application.NewPostService(store, createJobOfferOrchestrator)
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

func (server *Server) initCreateJobOfferOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.CreateJobOfferOrchestrator {
	orchestrator, err := application.NewCreateJobOfferOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
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

func (server *Server) initCreateJobOfferHandler(service *application.PostService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewJobOfferCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) startGrpcServer(postHandler *api.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	postProto.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
