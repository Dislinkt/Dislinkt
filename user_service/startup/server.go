package startup

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dislinkt/common/interceptor"
	userProto "github.com/dislinkt/common/proto/user_service"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/common/saga/messaging/nats"
	"github.com/dislinkt/common/tracer"
	"github.com/dislinkt/user_service/application"
	"github.com/dislinkt/user_service/domain"
	"github.com/dislinkt/user_service/infrastructure/api"
	"github.com/dislinkt/user_service/infrastructure/persistence"
	"github.com/dislinkt/user_service/startup/config"
	"github.com/gofrs/uuid"
	otgo "github.com/opentracing/opentracing-go"
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
	QueueGroupRegister = "user_service_register"
	QueueGroupUpdate   = "user_service_update"
	QueueGroupPatch    = "user_service_patch"
)

func (server *Server) Start() {

	tracer, _ := tracer.Init("user_service")
	otgo.SetGlobalTracer(tracer)

	postgresClient := server.initUserClient()
	userStore := server.initUserStore(postgresClient)
	server.addAdmin(userStore)

	commandPublisher := server.initPublisher(server.config.RegisterUserCommandSubject)
	replySubscriber := server.initSubscriber(server.config.RegisterUserReplySubject, QueueGroupRegister)
	registerUserOrchestrator := server.initRegisterUserOrchestrator(commandPublisher, replySubscriber)

	patchCommandPublisher := server.initPublisher(server.config.PatchUserCommandSubject)
	patchReplySubscriber := server.initSubscriber(server.config.PatchUserReplySubject, QueueGroupPatch)
	patchUserOrchestrator := server.initPatchUserOrchestrator(patchCommandPublisher, patchReplySubscriber)

	updateCommandPublisher := server.initPublisher(server.config.UpdateUserCommandSubject)
	updateReplySubscriber := server.initSubscriber(server.config.UpdateUserReplySubject, QueueGroupUpdate)
	updateUserOrchestrator := server.initUpdateUserOrchestrator(updateCommandPublisher, updateReplySubscriber)

	userService := server.initUserService(userStore, registerUserOrchestrator, updateUserOrchestrator, patchUserOrchestrator)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroupRegister)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(userService, replyPublisher, commandSubscriber)

	patchCommandSubscriber := server.initSubscriber(server.config.PatchUserCommandSubject, QueueGroupPatch)
	patchReplyPublisher := server.initPublisher(server.config.PatchUserReplySubject)
	server.initPatchUserHandler(userService, patchCommandSubscriber, patchReplyPublisher)

	updateCommandSubscriber := server.initSubscriber(server.config.UpdateUserCommandSubject, QueueGroupUpdate)
	updateReplyPublisher := server.initPublisher(server.config.UpdateUserReplySubject)
	server.initUpdateUserHandler(userService, updateReplyPublisher, updateCommandSubscriber)

	userHandler := server.initUserHandler(userService)
	server.startGrpcServer(userHandler)

}

func (server *Server) addAdmin(store domain.UserStore) {
	id, _ := uuid.FromString("13c4dc91-410c-4370-a964-17c64566f740")
	username := "admin"
	email := "admin@gmail.com"
	user := domain.User{
		Id:          id,
		Name:        "Admin",
		Surname:     "Admin",
		Username:    &username,
		Email:       &email,
		Number:      "",
		Gender:      0,
		DateOfBirth: "",
		Password:    "",
		UserRole:    1,
		Biography:   "",
		Blocked:     false,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Private:     true,
	}
	err := store.Insert(&user)
	if err != nil {
		log.Println("Error creating admin!")
		return
	}
	log.Println("Admin successfully added!")
}

func (server *Server) initUserClient() *gorm.DB {
	client, err := persistence.GetClient(
		server.config.UserDBHost, server.config.UserDBUser,
		server.config.UserDBPass, server.config.UserDBName,
		server.config.UserDBPort)
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
	// 	err := store.Register(Product)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	return store
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

func (server *Server) initRegisterUserOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.RegisterUserOrchestrator {
	orchestrator, err := application.NewRegisterUserOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initPatchUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *application.PatchUserOrchestrator {
	orchestrator, err := application.NewPatchUserOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}
func (server *Server) initUpdateUserOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.UpdateUserOrchestrator {
	orchestrator, err := application.NewUpdateUserOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initUserService(store domain.UserStore,
	registerUserOrchestrator *application.RegisterUserOrchestrator,
	updateUserOrchestrator *application.UpdateUserOrchestrator, patchOrchestrator *application.PatchUserOrchestrator) *application.UserService {
	return application.NewUserService(store, registerUserOrchestrator, updateUserOrchestrator, patchOrchestrator)
}

func (server *Server) initRegisterUserHandler(service *application.UserService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewRegisterUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initPatchUserHandler(service *application.UserService, subscriber saga.Subscriber, publisher saga.Publisher) {
	_, err := api.NewPatchUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}
func (server *Server) initUpdateUserHandler(service *application.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewUpdateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initUserHandler(service *application.UserService) *api.UserHandler {
	return api.NewUserHandler(service)
}

func (server *Server) startGrpcServer(userHandler *api.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	userProto.RegisterUserServiceServer(grpcServer, userHandler)
	// grpcServer := grpc.NewServer()
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
