package startup

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"path/filepath"

	logger "github.com/dislinkt/common/logging"

	"github.com/dislinkt/common/interceptor"
	"google.golang.org/grpc/credentials"

	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/common/saga/messaging/nats"
	"github.com/dislinkt/connection_service/infrastructure/persistance"
	"google.golang.org/grpc"

	connection "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/application"
	"github.com/dislinkt/connection_service/domain"
	"github.com/dislinkt/connection_service/infrastructure/api"
	"github.com/dislinkt/connection_service/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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
	QueueGroup = "connection_service"
)

func (server *Server) Start() {

	neo4jClient := server.initNeo4J()

	connectionStore := server.initConnectionStore(neo4jClient)

	connectionService := server.initConnectionService(connectionStore)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(connectionService, replyPublisher, commandSubscriber)

	patchCommandSubscriber := server.initSubscriber(server.config.PatchUserCommandSubject, QueueGroup)
	patchReplyPublisher := server.initPublisher(server.config.PatchUserReplySubject)
	server.iniPatchUserHandler(connectionService, patchReplyPublisher, patchCommandSubscriber)

	connectionHandler := server.initConnectionHandler(connectionService)

	server.startGrpcServer(connectionHandler)
	server.logger.InfoLogger.Info("SS")
}

func (server *Server) initNeo4J() *neo4j.Driver {
	client, err := persistance.GetClient(server.config.Neo4jUri, server.config.Neo4jUsername, server.config.Neo4jPassword)
	if err != nil {
		server.logger.ErrorLogger.Error("IC")
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
func (server *Server) initRegisterUserHandler(service *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewRegisterUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) iniPatchUserHandler(service *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewPatchUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}
func (server *Server) initConnectionHandler(service *application.ConnectionService) *api.ConnectionHandler {
	return api.NewConnectionHandler(service)
}

func (server *Server) startGrpcServer(connectionHandler *api.ConnectionHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.logger.ErrorLogger.Error("FTL")
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()), grpc.Creds(tlsCredentials))
	connection.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.logger.ErrorLogger.Error("FTS")
		log.Fatalf("failed to serve: %s", err)
	}
}
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	// caCert, _ := filepath.Abs("./ca-cert.pem")
	// pemClientCA, err := ioutil.ReadFile(caCert)
	// if err != nil {
	// 	return nil, err
	// }

	// certPool := x509.NewCertPool()
	// if !certPool.AppendCertsFromPEM(pemClientCA) {
	// 	return nil, fmt.Errorf("failed to add client CA's certificate")
	// }

	// Load server's certificate and private key
	crtPath, _ := filepath.Abs("./server-cert.pem")
	keyPath, _ := filepath.Abs("./server-key.pem")
	serverCert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
