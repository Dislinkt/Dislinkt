package startup

import (
	"fmt"
	"github.com/dislinkt/connection_service/infrastructure/persistance"
	"google.golang.org/grpc"

	connection "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/application"
	"github.com/dislinkt/connection_service/domain"
	"github.com/dislinkt/connection_service/infrastructure/api"
	"github.com/dislinkt/connection_service/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {

	neo4jClient := server.initNeo4J()

	connectionStore := server.initConnectionStore(neo4jClient)

	connectionService := server.initConnectionService(connectionStore)

	connectionHandler := server.initConnectionHandler(connectionService)

	server.startGrpcServer(connectionHandler)
}

func (server *Server) initNeo4J() *neo4j.Driver {
	client, err := persistance.GetClient(server.config.Neo4jUri, server.config.Neo4jUsername, server.config.Neo4jPassword)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initConnectionStore(client *neo4j.Driver) domain.ConnectionStore {
	store := persistance.NewConnectionDBStore(client)
	return store
}

func (server *Server) initConnectionService(store domain.ConnectionStore) *application.ConnectionService {
	return application.NewConnectionService(store)
}

func (server *Server) initConnectionHandler(service *application.ConnectionService) *api.ConnectionHandler {
	return api.NewConnectionHandler(service)
}

func (server *Server) startGrpcServer(connectionHandler *api.ConnectionHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	connection.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
