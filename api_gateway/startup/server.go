package startup

import (
	"context"
	"fmt"
	"github.com/dislinkt/api_gateway/infrastructure/api"
	authGw "github.com/dislinkt/common/proto/auth_service"
	"io"
	"log"
	"net/http"

	cfg "github.com/dislinkt/api_gateway/startup/config"
	additionaluserGw "github.com/dislinkt/common/proto/additional_user_service"
	connectionGw "github.com/dislinkt/common/proto/connection_service"
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
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initCustomHandlers()
	return server
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	err := userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts)
	additionalUserEndpoint := fmt.Sprintf("%s:%s", server.config.AdditionalUserHost, server.config.AdditionalUserPort)
	err = additionaluserGw.RegisterAdditionalUserServiceHandlerFromEndpoint(context.TODO(), server.mux, additionalUserEndpoint, opts)
	if err != nil {
		panic(err)
	}
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)
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
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	connectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	userFeedHandler := api.NewUserFeedHandler(postEndpoint, connectionEndpoint)
	userFeedHandler.Init(server.mux)
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux))
}
