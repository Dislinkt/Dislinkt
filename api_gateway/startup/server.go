package startup

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	logger "github.com/dislinkt/common/logging"

	"github.com/dislinkt/api_gateway/infrastructure/api"
	"github.com/gorilla/handlers"
	"google.golang.org/grpc/credentials"

	// "github.com/dislinkt/api_gateway/infrastructure/api"

	cfg "github.com/dislinkt/api_gateway/startup/config"
	additionalUserGw "github.com/dislinkt/common/proto/additional_user_service"
	authGw "github.com/dislinkt/common/proto/auth_service"
	connectionGw "github.com/dislinkt/common/proto/connection_service"
	postGw "github.com/dislinkt/common/proto/post_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
	logger *logger.Logger
	// tracer otgo.Tracer
	// closer io.Closer
}

func NewServer(config *cfg.Config) *Server {
	logger := logger.InitLogger(context.TODO())
	server := &Server{
		config: config,
		mux: runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(customMatcher),
		),
		logger: logger,
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
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCredentials)}
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	err = userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts)
	if err != nil {
		server.logger.ErrorLogger.Error("UMI")
		panic(err)
	}

	additionalUserEndpoint := fmt.Sprintf("%s:%s", server.config.AdditionalUserHost, server.config.AdditionalUserPort)
	err = additionalUserGw.RegisterAdditionalUserServiceHandlerFromEndpoint(context.TODO(), server.mux, additionalUserEndpoint, opts)
	if err != nil {
		server.logger.ErrorLogger.Error("AUMI")
		panic(err)
	}

	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)
	if err != nil {
		server.logger.ErrorLogger.Error("PMI")
		panic(err)
	}

	connectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionEndpoint, opts)
	if err != nil {
		server.logger.ErrorLogger.Error("CMI")
		panic(err)
	}

	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	err = authGw.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), server.mux, authEndpoint, opts)
	if err != nil {
		server.logger.ErrorLogger.Error("AMI")
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
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	caCert, _ := filepath.Abs("./ca-cert.pem")
	pemServerCA, err := ioutil.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}
	// Load client's certificate and private key
	// cCert, _ := filepath.Abs("./client-cert.pem")
	// clientKey, _ := filepath.Abs("./client-key.pem")
	// clientCert, err := tls.LoadX509KeyPair(cCert, clientKey)
	// if err != nil {
	// 	return nil, err
	// }

	// Create the credentials and return it
	config := &tls.Config{
		// Certificates: []tls.Certificate{clientCert},
		InsecureSkipVerify: true,
		RootCAs:            certPool,
	}

	return credentials.NewTLS(config), nil
}

func muxMiddleware(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(server.config.AuthHost + ":" + server.config.AuthPort)
		server.mux.ServeHTTP(w, r)
	})
}
