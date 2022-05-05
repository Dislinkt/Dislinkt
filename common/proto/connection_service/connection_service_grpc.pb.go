// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: connection_service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ConnectionServiceClient is the client API for ConnectionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectionServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	CreateConnection(ctx context.Context, in *NewConnectionRequest, opts ...grpc.CallOption) (*NewConnectionResponse, error)
}

type connectionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectionServiceClient(cc grpc.ClientConnInterface) ConnectionServiceClient {
	return &connectionServiceClient{cc}
}

func (c *connectionServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/connection_service_proto.ConnectionService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectionServiceClient) CreateConnection(ctx context.Context, in *NewConnectionRequest, opts ...grpc.CallOption) (*NewConnectionResponse, error) {
	out := new(NewConnectionResponse)
	err := c.cc.Invoke(ctx, "/connection_service_proto.ConnectionService/CreateConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectionServiceServer is the server API for ConnectionService service.
// All implementations must embed UnimplementedConnectionServiceServer
// for forward compatibility
type ConnectionServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	CreateConnection(context.Context, *NewConnectionRequest) (*NewConnectionResponse, error)
	mustEmbedUnimplementedConnectionServiceServer()
}

// UnimplementedConnectionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConnectionServiceServer struct {
}

func (UnimplementedConnectionServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedConnectionServiceServer) CreateConnection(context.Context, *NewConnectionRequest) (*NewConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateConnection not implemented")
}
func (UnimplementedConnectionServiceServer) mustEmbedUnimplementedConnectionServiceServer() {}

// UnsafeConnectionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectionServiceServer will
// result in compilation errors.
type UnsafeConnectionServiceServer interface {
	mustEmbedUnimplementedConnectionServiceServer()
}

func RegisterConnectionServiceServer(s grpc.ServiceRegistrar, srv ConnectionServiceServer) {
	s.RegisterService(&ConnectionService_ServiceDesc, srv)
}

func _ConnectionService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/connection_service_proto.ConnectionService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectionService_CreateConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServiceServer).CreateConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/connection_service_proto.ConnectionService/CreateConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServiceServer).CreateConnection(ctx, req.(*NewConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConnectionService_ServiceDesc is the grpc.ServiceDesc for ConnectionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConnectionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "connection_service_proto.ConnectionService",
	HandlerType: (*ConnectionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _ConnectionService_Register_Handler,
		},
		{
			MethodName: "CreateConnection",
			Handler:    _ConnectionService_CreateConnection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "connection_service.proto",
}
