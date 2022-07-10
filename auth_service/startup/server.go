package startup

import (
	"errors"
	"fmt"
	"github.com/dislinkt/common/tracer"
	otgo "github.com/opentracing/opentracing-go"
	"log"
	"net"
	"regexp"

	"github.com/dislinkt/common/interceptor"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/common/saga/messaging/nats"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

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
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

const (
	QueueGroup       = "auth_service"
	QueueGroupUpdate = "auth_service_update"
)

func (server *Server) Start() {

	tracer, _ := tracer.Init("auth_service")
	otgo.SetGlobalTracer(tracer)

	postgresClient := server.initUserClient()
	userStore := server.initUserStore(postgresClient)
	permissionStore := server.initPermissionStore(postgresClient)
	server.addAdmin(userStore)

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

}

func (server *Server) initUserClient() *gorm.DB {
	client, err := persistence.GetClient(
		server.config.AuthDBHost, server.config.AuthDBUser,
		server.config.AuthDBPass, server.config.AuthDBName,
		server.config.AuthDBPort)
	if err != nil {
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
		log.Fatalf("failed to listen: %v", err)
	}

	// publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	// if err != nil {
	//	log.Fatalf("failed to parse public key: %v", err)
	// }

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	authProto.RegisterAuthServiceServer(grpcServer, authHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initUpdateUserHandler(service *application.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewUpdateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) addAdmin(store domain.UserStore) {
	id, _ := uuid.FromString("13c4dc91-410c-4370-a964-17c64566f740")
	hashAndSalt, err := HashAndSaltPasswordIfStrongAndMatching("Lozinka123")
	user := domain.User{
		Id:        id,
		Username:  "admin",
		Email:     "admin@gmail.com",
		Password:  hashAndSalt,
		UserRole:  1,
		Active:    true,
		ApiToken:  nil,
		TotpToken: "",
	}
	err = store.Insert(&user)
	if err != nil {
		log.Println("Error creating admin!")
		return
	}
	log.Println("Admin successfully added!")
}

func HashAndSaltPasswordIfStrongAndMatching(password string) (string, error) {
	isStrong, _ := regexp.MatchString("[0-9A-Za-z!?#$@.*+_\\-]+", password)

	if !isStrong {
		return "", errors.New("Password not strong enough!")
	}
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash), err
}
