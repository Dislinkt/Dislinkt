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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/credentials"

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
	logger *logger.Logger
	// tracer otgo.Tracer
	// closer io.Closer
}

func NewServer(config *config.Config) *Server {
	logger := logger.InitLogger(context.TODO())
	// newTracer, closer := tracer.Init(config.JaegerServiceName)
	// otgo.SetGlobalTracer(newTracer)
	return &Server{
		config: config,
		logger: logger,
		// tracer: newTracer,
		// closer: closer,
	}
}

const (
	QueueGroup = "additional_user_service"
)

func (server *Server) Start() {
	mongoClient := server.initAdditionalUserClient()
	additionalUserStore := server.initAdditionalUserStore(mongoClient)

	additionalUserService := server.initAdditionalUserService(additionalUserStore)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(additionalUserService, replyPublisher, commandSubscriber)

	additionalUserHandler := server.initAdditionalUserHandler(additionalUserService)

	server.initData(additionalUserService, additionalUserStore)

	server.startGrpcServer(additionalUserHandler)
	server.logger.InfoLogger.Info("SS")
}

func (server *Server) initAdditionalUserClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.AdditionalUserDBHost, server.config.AdditionalUserDBPort)
	if err != nil {
		server.logger.ErrorLogger.Error("IC")
		log.Fatal(err)
	}
	return client
}

func (server *Server) initAdditionalUserStore(client *mongo.Client) domain.AdditionalUserStore {
	store := persistence.NewAdditionalUserMongoDBStore(client)
	return store
}

func (server *Server) initAdditionalUserService(store domain.AdditionalUserStore) *application.AdditionalUserService {
	return application.NewAdditionalUserService(store)
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

func (server *Server) initAdditionalUserHandler(service *application.AdditionalUserService) *api.AdditionalUserHandler {
	return api.NewProductHandler(service)
}

func (server *Server) startGrpcServer(additionalUserHandler *api.AdditionalUserHandler) {
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
	additionalUserProto.RegisterAdditionalUserServiceServer(grpcServer, additionalUserHandler)
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

func (server *Server) initData(service *application.AdditionalUserService, store domain.AdditionalUserStore) {
	var fields []*domain.FieldOfStudy
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Business Administration and Management, General"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Electrical and Electronics Engineering"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Accounting"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "English Language and Literature/Letters"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Political Science and Government"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Computer and Information Sciences and Support Services"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Communication and Media Studies"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Human Resources Management/Personnel Administration/General"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Architecture"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Electronic, Electronics and Communications Engineering"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Liberal Arts and Sciences/Liberal Studies"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(),
		Name: "International Relations and Affairs"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Computer Systems networking and Telecommunications"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Criminal Justice and Corrections"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Business, Management, Marketing, and Related Support Services"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Art/Art Studies, General"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Advertising"})
	fields = append(fields, &domain.FieldOfStudy{Id: primitive.NewObjectID(), Name: "Fine/Studio Arts, General"})
	_, err := store.InsertFieldOfStudy(fields)
	if err != nil {
		server.logger.WarnLogger.Warn("FOSC")
		return
	}
	server.logger.InfoLogger.Info("FOSC")

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
		server.logger.WarnLogger.Warn("SC")
		return
	}
	server.logger.InfoLogger.Info("SC")

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
		server.logger.WarnLogger.Warn("IC")
		return
	}
	server.logger.InfoLogger.Info("IC")

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
		server.logger.WarnLogger.Warn("DC")
		return
	}
	server.logger.InfoLogger.Info("DC")

}
