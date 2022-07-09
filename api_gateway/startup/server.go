package startup

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dislinkt/api_gateway/infrastructure/api"
	"github.com/gorilla/handlers"

	// "github.com/dislinkt/api_gateway/infrastructure/api"

	cfg "github.com/dislinkt/api_gateway/startup/config"
	additionalUserGw "github.com/dislinkt/common/proto/additional_user_service"
	authGw "github.com/dislinkt/common/proto/auth_service"
	connectionGw "github.com/dislinkt/common/proto/connection_service"
	eventGw "github.com/dislinkt/common/proto/event_service"
	messageGw "github.com/dislinkt/common/proto/message_service"
	notificationGw "github.com/dislinkt/common/proto/notification_service"
	postGw "github.com/dislinkt/common/proto/post_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	otgo "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
	tracer otgo.Tracer
	closer io.Closer
}

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux: runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(customMatcher),
		),
	}
	server.initHandlers()
	server.initCustomHandlers()
	return server
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Authorization":
		return key, true
	default:
		return key, false
	}
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	err := userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts)
	if err != nil {
		panic(err)
	}

	additionalUserEndpoint := fmt.Sprintf("%s:%s", server.config.AdditionalUserHost, server.config.AdditionalUserPort)
	err = additionalUserGw.RegisterAdditionalUserServiceHandlerFromEndpoint(context.TODO(), server.mux, additionalUserEndpoint, opts)
	if err != nil {
		panic(err)
	}

	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)
	if err != nil {
		panic(err)
	}

	messageEndpoint := fmt.Sprintf("%s:%s", server.config.MessageHost, server.config.MessagePort)
	err = messageGw.RegisterMessageServiceHandlerFromEndpoint(context.TODO(), server.mux, messageEndpoint, opts)
	if err != nil {
		panic(err)
	}

	notificationEndpoint := fmt.Sprintf("%s:%s", server.config.NotificationHost, server.config.NotificationPort)
	err = notificationGw.RegisterNotificationServiceHandlerFromEndpoint(context.TODO(), server.mux, notificationEndpoint, opts)
	if err != nil {
		panic(err)
	}

	eventEndpoint := fmt.Sprintf("%s:%s", server.config.EventHost, server.config.EventPort)
	err = eventGw.RegisterEventServiceHandlerFromEndpoint(context.TODO(), server.mux, eventEndpoint, opts)
	if err != nil {
		panic(err)
	}

	connectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionEndpoint, opts)
	if err != nil {
		panic(err)
	}

	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	err = authGw.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), server.mux, authEndpoint, opts)
	if err != nil {
		panic(err)
	}
}

func (server *Server) initCustomHandlers() {
	// postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	// connectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	userFeedHandler := api.NewUserFeedHandler(server.config)
	userFeedHandler.Init(server.mux)
	connectionRequestHandler := api.NewConnectionRequestHandler(server.config)
	connectionRequestHandler.Init(server.mux)
	connectionBlockHandler := api.NewConnectionBlockedUsersHandler(server.config)
	connectionBlockHandler.Init(server.mux)
	connectionAllHandler := api.NewConnectionAllHandler(server.config)
	connectionAllHandler.Init(server.mux)
	publicUserFeedHandler := api.NewPublicUserFeedHandler(server.config)
	publicUserFeedHandler.Init(server.mux)
	recommendationsHandler := api.NewConnectionRecommendationHandler(server.config)
	recommendationsHandler.Init(server.mux)
	searchUserHandler := api.NewConnectionSearchUserHandler(server.config)
	searchUserHandler.Init(server.mux)
}

func (server *Server) Start() {
	crtPath, _ := filepath.Abs("./cert.crt")
	keyPath, _ := filepath.Abs("./cert.key")
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"https://localhost:4200", "https://localhost:4200/**",
			"http://localhost:4200", "http://localhost:4200/**", "http://localhost:8080/**",
			"http://localhost:3000/**", "http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "Authorization", "Access-Control-Allow-Origin", "*"}),
		handlers.AllowCredentials(),
	)
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), cors(muxMiddleware(server))))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", server.config.Port), crtPath, keyPath, cors(muxMiddleware(server))))
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux))
}

func muxMiddleware(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(server.config.AuthHost + ":" + server.config.AuthPort)
		server.mux.ServeHTTP(w, r)
	})
}
