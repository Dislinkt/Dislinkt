package startup

import (
	"fmt"
	"github.com/dislinkt/common/tracer"
	otgo "github.com/opentracing/opentracing-go"
	"log"
	"net"

	"github.com/dislinkt/common/interceptor"

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
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

const (
	QueueGroup = "connection_service"
)

func (server *Server) Start() {
	tracer, _ := tracer.Init("user_service")
	otgo.SetGlobalTracer(tracer)

	neo4jClient := server.initNeo4J()

	connectionStore := server.initConnectionStore(neo4jClient)

	server.initData(connectionStore)
	connectionService := server.initConnectionService(connectionStore)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(connectionService, replyPublisher, commandSubscriber)

	patchCommandSubscriber := server.initSubscriber(server.config.PatchUserCommandSubject, QueueGroup)
	patchReplyPublisher := server.initPublisher(server.config.PatchUserReplySubject)
	server.iniPatchUserHandler(connectionService, patchReplyPublisher, patchCommandSubscriber)

	fmt.Println(server.config.CreateJobOfferCommandSubject + " start connection_service")
	createJobOfferSubscriber := server.initSubscriber(server.config.CreateJobOfferCommandSubject, QueueGroup)
	createJobOfferReplyPublisher := server.initPublisher(server.config.CreateJobOfferReplySubject)
	server.initCreateJobOfferHandler(connectionService, createJobOfferReplyPublisher, createJobOfferSubscriber)

	commandAddEducationSubscriber := server.initSubscriber(server.config.AddEducationCommandSubject, QueueGroup)
	replyAddEducationPublisher := server.initPublisher(server.config.AddEducationReplySubject)
	server.initAddEducationHandler(connectionService, replyAddEducationPublisher, commandAddEducationSubscriber)

	commandAddSkillSubscriber := server.initSubscriber(server.config.AddSkillCommandSubject, QueueGroup)
	replyAddSkillPublisher := server.initPublisher(server.config.AddSkillReplySubject)
	server.initAddSkillHandler(connectionService, replyAddSkillPublisher, commandAddSkillSubscriber)

	commandDeleteSkillSubscriber := server.initSubscriber(server.config.DeleteSkillCommandSubject, QueueGroup)
	replyDeleteSkillPublisher := server.initPublisher(server.config.DeleteSkillReplySubject)
	server.initDeleteSkillHandler(connectionService, replyDeleteSkillPublisher, commandDeleteSkillSubscriber)

	commandUpdateSkillSubscriber := server.initSubscriber(server.config.UpdateSkillCommandSubject, QueueGroup)
	replyUpdateSkillPublisher := server.initPublisher(server.config.UpdateSkillReplySubject)
	server.initUpdateSkillHandler(connectionService, replyUpdateSkillPublisher, commandUpdateSkillSubscriber)

	commandDeleteEducationSubscriber := server.initSubscriber(server.config.DeleteEducationCommandSubject, QueueGroup)
	replyDeleteEducationPublisher := server.initPublisher(server.config.DeleteEducationReplySubject)
	server.initDeleteEducationHandler(connectionService, replyDeleteEducationPublisher, commandDeleteEducationSubscriber)

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

func (server *Server) initCreateJobOfferHandler(service *application.ConnectionService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewJobOfferCommandHandler(service, publisher, subscriber)
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

func (server *Server) initData(store domain.ConnectionStore) {
	isFilled, err := store.CheckIfDatabaseFilled()
	fmt.Println(isFilled)
	if !isFilled {

		_, err := store.DeleteAllFields()
		if err != nil {
			return
		}
		_, err = store.InsertField("Business Administration and Management, General")
		if err != nil {
			return
		}
		_, err = store.InsertField("Electrical and Electronics Engineering")
		if err != nil {
			return
		}
		_, err = store.InsertField("Accounting")
		if err != nil {
			return
		}
		_, err = store.InsertField("English Language and Literature/Letters")
		if err != nil {
			return
		}
		_, err = store.InsertField("Political Science and Government")
		if err != nil {
			return
		}
		_, err = store.InsertField("Computer and Information Sciences and Support Services")
		if err != nil {
			return
		}
		_, err = store.InsertField("Communication and Media Studies")
		if err != nil {
			return
		}

		_, err = store.DeleteAllSkills()
		_, err = store.InsertSkill("Communication")
		_, err = store.InsertSkill("Teamwork")
		_, err = store.InsertSkill("Critical Thinking")
		_, err = store.InsertSkill("Active Listening")
		_, err = store.InsertSkill("Active Learning")
		_, err = store.InsertSkill("Problem Solving")
		_, err = store.InsertSkill("Management")
		_, err = store.InsertSkill("Training")
		_, err = store.InsertSkill("Design")
		_, err = store.InsertSkill("Presentations")
		_, err = store.InsertSkill("Data Analysis")
		_, err = store.InsertSkill("Blogging")
		_, err = store.InsertSkill("Business")
		_, err = store.InsertSkill("Leadership")
		_, err = store.InsertSkill("Time Management")
		_, err = store.InsertSkill("Troubleshooting")
		_, err = store.InsertSkill("Operating System")
		_, err = store.InsertSkill("Online Marketing")
	}

	fmt.Println("zavrsiloo")

	if err != nil {
		return
	}

}

func (server *Server) startGrpcServer(connectionHandler *api.ConnectionHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	connection.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initAddEducationHandler(service *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewAddEducationCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initAddSkillHandler(service *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewAddSkillCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initDeleteSkillHandler(service *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewDeleteSkillCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initDeleteEducationHandler(service *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewDeleteEducationCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initUpdateSkillHandler(service *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewUpdateSkillCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}
