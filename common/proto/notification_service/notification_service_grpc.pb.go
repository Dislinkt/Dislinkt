// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: notification_service.proto

package notification_service

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

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	GetNotificationsForUser(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetMultipleResponse, error)
	SaveNotification(ctx context.Context, in *SaveNotificationRequest, opts ...grpc.CallOption) (*Empty, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) GetNotificationsForUser(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetMultipleResponse, error) {
	out := new(GetMultipleResponse)
	err := c.cc.Invoke(ctx, "/notification_service_proto.NotificationService/getNotificationsForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) SaveNotification(ctx context.Context, in *SaveNotificationRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/notification_service_proto.NotificationService/saveNotification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility
type NotificationServiceServer interface {
	GetNotificationsForUser(context.Context, *Empty) (*GetMultipleResponse, error)
	SaveNotification(context.Context, *SaveNotificationRequest) (*Empty, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNotificationServiceServer struct {
}

func (UnimplementedNotificationServiceServer) GetNotificationsForUser(context.Context, *Empty) (*GetMultipleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotificationsForUser not implemented")
}
func (UnimplementedNotificationServiceServer) SaveNotification(context.Context, *SaveNotificationRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveNotification not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_GetNotificationsForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).GetNotificationsForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notification_service_proto.NotificationService/getNotificationsForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).GetNotificationsForUser(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_SaveNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SaveNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notification_service_proto.NotificationService/saveNotification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SaveNotification(ctx, req.(*SaveNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification_service_proto.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getNotificationsForUser",
			Handler:    _NotificationService_GetNotificationsForUser_Handler,
		},
		{
			MethodName: "saveNotification",
			Handler:    _NotificationService_SaveNotification_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notification_service.proto",
}
