module github.com/dislinkt/connection_service

go 1.18

replace github.com/dislinkt/common => ../common

require (
	github.com/dislinkt/common v0.0.0-00010101000000-000000000000
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/neo4j/neo4j-go-driver/v4 v4.4.2
	github.com/opentracing/opentracing-go v1.2.0
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/nats-io/nats-server/v2 v2.8.2 // indirect
	github.com/nats-io/nats.go v1.15.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.mongodb.org/mongo-driver v1.9.1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220426171045-31bebdecfb46 // indirect
)
