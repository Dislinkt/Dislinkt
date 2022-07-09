package startup

import (
	"fmt"
	"github.com/dislinkt/common/interceptor"
	eventProto "github.com/dislinkt/common/proto/event_service"
	"github.com/dislinkt/common/tracer"
	"github.com/dislinkt/event_service/application"
	"github.com/dislinkt/event_service/domain"
	"github.com/dislinkt/event_service/infrastructure/api"
	"github.com/dislinkt/event_service/infrastructure/persistence"
	"github.com/dislinkt/event_service/startup/config"
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
	tracer, _ := tracer.Init("user_service")
	otgo.SetGlobalTracer(tracer)
	mongoClient := server.initMongoClient()
	eventStore := server.initEventStore(mongoClient)

	eventService := server.initEventService(eventStore)
	eventHandler := server.initEventHandler(eventService)
	server.startGrpcServer(eventHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.EventDBHost, server.config.EventDBPort)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func (server *Server) initEventStore(client *mongo.Client) domain.EventStore {
	store := persistence.NewEventMongoDBStore(client)
	return store
}

func (server *Server) initEventService(store domain.EventStore) *application.EventService {
	return application.NewEventService(store)
}

func (server *Server) initEventHandler(service *application.EventService) *api.EventHandler {
	return api.NewEventHandler(service)
}

func (server *Server) startGrpcServer(eventHandler *api.EventHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	eventProto.RegisterEventServiceServer(grpcServer, eventHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
