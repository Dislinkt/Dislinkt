package startup

import (
	"fmt"
	"github.com/dislinkt/common/tracer"
	otgo "github.com/opentracing/opentracing-go"
	"log"
	"net"

	"github.com/dislinkt/common/interceptor"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dislinkt/additional_user_service/application"
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/additional_user_service/infrastructure/api"
	"github.com/dislinkt/additional_user_service/infrastructure/persistence"
	"github.com/dislinkt/additional_user_service/startup/config"
	additionalUserProto "github.com/dislinkt/common/proto/additional_user_service"
	saga "github.com/dislinkt/common/saga/messaging"
	"github.com/dislinkt/common/saga/messaging/nats"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
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
	QueueGroup = "additional_user_service"
)

func (server *Server) Start() {
	tracer, _ := tracer.Init("additional_user_service")
	otgo.SetGlobalTracer(tracer)
	mongoClient := server.initAdditionalUserClient()
	additionalUserStore := server.initAdditionalUserStore(mongoClient)

	commandPublisher := server.initPublisher(server.config.AddEducationCommandSubject)
	replySubscriber := server.initSubscriber(server.config.AddEducationReplySubject, QueueGroup)
	addEducationOrchestrator := server.initAddEducationOrchestrator(commandPublisher, replySubscriber)

	commandSkillPublisher := server.initPublisher(server.config.AddSkillCommandSubject)
	replySkillSubscriber := server.initSubscriber(server.config.AddSkillReplySubject, QueueGroup)
	addSkillOrchestrator := server.initSkillOrchestrator(commandSkillPublisher, replySkillSubscriber)

	commandDeleteSkillPublisher := server.initPublisher(server.config.DeleteSkillCommandSubject)
	replyDeleteSkillSubscriber := server.initSubscriber(server.config.DeleteSkillReplySubject, QueueGroup)
	deleteSkillOrchestrator := server.initDeleteSkillOrchestrator(commandDeleteSkillPublisher, replyDeleteSkillSubscriber)

	commandUpdateSkillPublisher := server.initPublisher(server.config.UpdateSkillCommandSubject)
	replyUpdateSkillSubscriber := server.initSubscriber(server.config.UpdateSkillReplySubject, QueueGroup)
	updateSkillOrchestrator := server.initUpdateSkillOrchestrator(commandUpdateSkillPublisher, replyUpdateSkillSubscriber)

	commandDeleteEducationPublisher := server.initPublisher(server.config.DeleteEducationCommandSubject)
	replyDeleteEducationSubscriber := server.initSubscriber(server.config.DeleteEducationReplySubject, QueueGroup)
	deleteEducationOrchestrator := server.initDeleteEducationOrchestrator(commandDeleteEducationPublisher, replyDeleteEducationSubscriber)

	commandUpdateEducationPublisher := server.initPublisher(server.config.UpdateEducationCommandSubject)
	replyUpdateEducationSubscriber := server.initSubscriber(server.config.UpdateEducationReplySubject, QueueGroup)
	updateEducationOrchestrator := server.initUpdateEducationOrchestrator(commandUpdateEducationPublisher, replyUpdateEducationSubscriber)

	additionalUserService := server.initAdditionalUserService(additionalUserStore, addEducationOrchestrator, addSkillOrchestrator, deleteSkillOrchestrator, updateSkillOrchestrator,
		deleteEducationOrchestrator, updateEducationOrchestrator)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(additionalUserService, replyPublisher, commandSubscriber)

	commandAddEducationSubscriber := server.initSubscriber(server.config.AddEducationCommandSubject, QueueGroup)
	replyAddEducationPublisher := server.initPublisher(server.config.AddEducationReplySubject)
	server.initAddEducationHandler(additionalUserService, replyAddEducationPublisher, commandAddEducationSubscriber)

	commandAddSkillSubscriber := server.initSubscriber(server.config.AddSkillCommandSubject, QueueGroup)
	replyAddSkillPublisher := server.initPublisher(server.config.AddSkillReplySubject)
	server.initAddSkillHandler(additionalUserService, replyAddSkillPublisher, commandAddSkillSubscriber)

	commandeleteSkillSubscriber := server.initSubscriber(server.config.DeleteSkillCommandSubject, QueueGroup)
	replyDeleteSkillPublisher := server.initPublisher(server.config.DeleteSkillReplySubject)
	server.initDeleteSkillHandler(additionalUserService, replyDeleteSkillPublisher, commandeleteSkillSubscriber)

	commandUpdateSkillSubscriber := server.initSubscriber(server.config.UpdateSkillCommandSubject, QueueGroup)
	replyUpdateSkillPublisher := server.initPublisher(server.config.UpdateSkillReplySubject)
	server.initUpdateSkillHandler(additionalUserService, replyUpdateSkillPublisher, commandUpdateSkillSubscriber)

	commandDeleteEducationSubscriber := server.initSubscriber(server.config.DeleteEducationCommandSubject, QueueGroup)
	replyDeleteEducationPublisher := server.initPublisher(server.config.DeleteEducationReplySubject)
	server.initDeleteEducationHandler(additionalUserService, replyDeleteEducationPublisher, commandDeleteEducationSubscriber)

	commandUpdateEducationSubscriber := server.initSubscriber(server.config.UpdateEducationCommandSubject, QueueGroup)
	replyUpdateEducationPublisher := server.initPublisher(server.config.UpdateEducationReplySubject)
	server.initUpdateEducationHandler(additionalUserService, replyUpdateEducationPublisher, commandUpdateEducationSubscriber)

	additionalUserHandler := server.initAdditionalUserHandler(additionalUserService)

	server.initData(additionalUserService, additionalUserStore)

	server.startGrpcServer(additionalUserHandler)
}

func (server *Server) initAdditionalUserClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.AdditionalUserDBHost, server.config.AdditionalUserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initAdditionalUserStore(client *mongo.Client) domain.AdditionalUserStore {
	store := persistence.NewAdditionalUserMongoDBStore(client)
	return store
}

func (server *Server) initAdditionalUserService(store domain.AdditionalUserStore, addEducationOrchestrator *application.AddEducationOrchestrator,
	addSkillOrchestrator *application.AddSkillOrchestrator, deleteSkillOrchestrator *application.DeleteSkillOrchestrator,
	updateSkillOrchestrator *application.UpdateSkillOrchestrator, deleteEducationOrchestrator *application.DeleteEducationOrchestrator,
	updateEducationOrchestrator *application.UpdateEducationOrchestrator) *application.AdditionalUserService {
	return application.NewAdditionalUserService(store, addEducationOrchestrator, addSkillOrchestrator, deleteSkillOrchestrator, updateSkillOrchestrator, deleteEducationOrchestrator,
		updateEducationOrchestrator)
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

func (server *Server) initRegisterUserHandler(service *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewRegisterUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initAddEducationHandler(service *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewAddEducationCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initAddSkillHandler(service *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewAddSkillCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initDeleteSkillHandler(service *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewDeleteSkillCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initDeleteEducationHandler(service *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewDeleteEducationCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initUpdateSkillHandler(service *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewUpdateSkillCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initUpdateEducationHandler(service *application.AdditionalUserService, publisher saga.Publisher,
	subscriber saga.Subscriber) {
	_, err := api.NewUpdateEducationCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initAdditionalUserHandler(service *application.AdditionalUserService) *api.AdditionalUserHandler {
	return api.NewProductHandler(service)
}

func (server *Server) startGrpcServer(additionalUserHandler *api.AdditionalUserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), server.config.PublicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	additionalUserProto.RegisterAdditionalUserServiceServer(grpcServer, additionalUserHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initData(service *application.AdditionalUserService, store domain.AdditionalUserStore) {
	var fields []*domain.FieldOfStudy
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Business Administration and Management, General"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Electrical and Electronics Engineering"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Accounting"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "English Language and Literature/Letters"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Political Science and Government"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Computer and Information Sciences and Support Services"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Communication and Media Studies"})
	_, err := store.InsertFieldOfStudy(fields)
	if err != nil {
		return
	}

	var skills []*domain.Skill
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Communication"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Teamwork"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Critical Thinking"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Active Listening"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Active Learning"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Problem Solving"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Management"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Training"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Design"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Presentations"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Data Analysis"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Blogging"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Business"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Leadership"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Time Management"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Troubleshooting"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Operating System"})
	skills = append(skills, &domain.Skill{Id: primitive.NewObjectID(), Name: "Online Marketing"})
	_, err = store.InsertSkills(skills)
	if err != nil {
		return
	}

	var industries []*domain.Industry
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(),
		Name: "IT Services and IT Consulting"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Hospitals and Health Care"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(),
		Name: "Education Administration Programs"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Government Administration"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Advertising Services"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Accounting"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Oil and Gas"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(),
		Name: "Wellness and Fitness Services"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Food and Beverage Services"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Appliances, Electrical, " +
		"and Electronics Manufacturing"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(),
		Name: "Business Consulting and Services"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Transportation, " +
		"Logistics and Storage"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Retail Apparel and Fashion"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(),
		Name: "Food and Beverage Manufacturing"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Staffing and Recruiting"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Architecture and Planning"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Travel Arrangements"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Armed Forces"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Airlines and Aviation"})
	industries = append(industries, &domain.Industry{Id: primitive.NewObjectID(), Name: "Blogs"})
	_, err = store.InsertIndustries(industries)
	if err != nil {
		return
	}

	var degrees []*domain.Degree
	degrees = append(degrees, &domain.Degree{Id: primitive.NewObjectID(),
		Name: "Doctoral degree"})
	degrees = append(degrees, &domain.Degree{Id: primitive.NewObjectID(), Name: "Professional degree"})
	degrees = append(degrees, &domain.Degree{Id: primitive.NewObjectID(),
		Name: "Master's degree"})
	degrees = append(degrees, &domain.Degree{Id: primitive.NewObjectID(), Name: "Bachelor's degree"})
	degrees = append(degrees, &domain.Degree{Id: primitive.NewObjectID(), Name: "Associate's degree"})
	degrees = append(degrees, &domain.Degree{Id: primitive.NewObjectID(), Name: "High school diploma"})
	_, err = store.InsertDegrees(degrees)
	if err != nil {
		return
	}

}

func (server *Server) initAddEducationOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.AddEducationOrchestrator {
	orchestrator, err := application.NewAddEducationOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initSkillOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.AddSkillOrchestrator {
	orchestrator, err := application.NewAddSkillOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initDeleteSkillOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.DeleteSkillOrchestrator {
	orchestrator, err := application.NewDeleteSkillOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initDeleteEducationOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.DeleteEducationOrchestrator {
	orchestrator, err := application.NewDeleteEducationOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initUpdateSkillOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.UpdateSkillOrchestrator {
	orchestrator, err := application.NewUpdateSkillOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initUpdateEducationOrchestrator(publisher saga.Publisher,
	subscriber saga.Subscriber) *application.UpdateEducationOrchestrator {
	orchestrator, err := application.NewUpdateEducationOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}
