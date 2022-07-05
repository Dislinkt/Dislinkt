module github.com/dislinkt/api_gateway

go 1.17

replace github.com/dislinkt/common => ../common

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0
	github.com/opentracing/opentracing-go v1.2.0
	google.golang.org/grpc v1.46.0
)

require (
	github.com/dislinkt/common v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.5.1
	github.com/prometheus/client_golang v1.12.2
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220426171045-31bebdecfb46 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
