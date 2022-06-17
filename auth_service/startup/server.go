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
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/common/saga/messaging/nats"
	"google.golang.org/grpc/credentials"

	"github.com/dislinkt/auth_service/application"
	"github.com/dislinkt/auth_service/domain"
	"github.com/dislinkt/auth_service/infrastructure/api"
	"github.com/dislinkt/auth_service/infrastructure/persistence"
	"github.com/dislinkt/auth_service/startup/config"
	authProto "github.com/dislinkt/common/proto/auth_service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
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
	QueueGroup       = "auth_service"
	QueueGroupUpdate = "auth_service_update"
)

func (server *Server) Start() {
	postgresClient := server.initUserClient()
	userStore := server.initUserStore(postgresClient)
	permissionStore := server.initPermissionStore(postgresClient)

	userService := server.initUserService(userStore)

	authService := server.initAuthService(userService, permissionStore)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(userService, authService, replyPublisher, commandSubscriber)

	updateCommandSubscriber := server.initSubscriber(server.config.UpdateUserCommandSubject, QueueGroupUpdate)
	updateReplyPublisher := server.initPublisher(server.config.UpdateUserReplySubject)
	server.initUpdateUserHandler(userService, updateReplyPublisher, updateCommandSubscriber)

	authHandler := server.initAuthHandler(authService)

	server.startGrpcServer(authHandler)
	server.logger.InfoLogger.Infof("SS")
}

func (server *Server) initUserClient() *gorm.DB {
	client, err := persistence.GetClient(
		server.config.AuthDBHost, server.config.AuthDBUser,
		server.config.AuthDBPass, server.config.AuthDBName,
		server.config.AuthDBPort)
	if err != nil {
		server.logger.ErrorLogger.Error("IC")
		log.Fatal(err)
	}
	return client
}

func (server *Server) initUserStore(client *gorm.DB) domain.UserStore {
	store, err := persistence.NewUserPostgresStore(client)
	if err != nil {
		log.Fatal(err)
	}
	// store.DeleteAll()
	// for _, Product := range products {
	// 	err := store.Insert(Product)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	return store
}

func (server *Server) initPermissionStore(client *gorm.DB) domain.PermissionStore {
	store, err := persistence.NewPermissionPostgresStore(client)
	if err != nil {
		log.Fatal(err)
	}
	err = store.DeleteAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, permission := range permissions {
		err := store.Insert(permission)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initUserService(store domain.UserStore) *application.UserService {
	return application.NewUserService(store)
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

func (server *Server) initRegisterUserHandler(service *application.UserService, authService *application.AuthService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewRegisterUserCommandHandler(service, authService, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initAuthService(userService *application.UserService, store domain.PermissionStore) *application.AuthService {
	return application.NewAuthService(userService, store)
}

func (server *Server) initAuthHandler(service *application.AuthService) *api.AuthHandler {
	return api.NewAuthHandler(service)
}

func (server *Server) startGrpcServer(authHandler *api.AuthHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.logger.ErrorLogger.Error("FTL")
		log.Fatalf("failed to listen: %v", err)
	}

	// publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	// if err != nil {
	//	log.Fatalf("failed to parse public key: %v", err)
	// }

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()), grpc.Creds(tlsCredentials))
	authProto.RegisterAuthServiceServer(grpcServer, authHandler)
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

func (server *Server) initUpdateUserHandler(service *application.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewUpdateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}
