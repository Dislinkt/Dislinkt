// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: user_service.proto

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

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	GetAll(ctx context.Context, in *SearchMessage, opts ...grpc.CallOption) (*GetAllResponse, error)
	GetOne(ctx context.Context, in *GetOneMessage, opts ...grpc.CallOption) (*UserResponse, error)
	GetMe(ctx context.Context, in *GetMeMessage, opts ...grpc.CallOption) (*GetMeResponse, error)
	RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	PatchUser(ctx context.Context, in *PatchUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	GetUserByUsername(ctx context.Context, in *GetOneByUsernameMessage, opts ...grpc.CallOption) (*UserResponse, error)
	GetPublicUsers(ctx context.Context, in *GetMeMessage, opts ...grpc.CallOption) (*GetAllResponse, error)
	CheckIfUserIsPrivate(ctx context.Context, in *GetOneMessage, opts ...grpc.CallOption) (*IsPrivateResponse, error)
	GetNotificationSettings(ctx context.Context, in *GetOneMessage, opts ...grpc.CallOption) (*NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, in *UpdateNotificationSettingsRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) GetAll(ctx context.Context, in *SearchMessage, opts ...grpc.CallOption) (*GetAllResponse, error) {
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetOne(ctx context.Context, in *GetOneMessage, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetMe(ctx context.Context, in *GetMeMessage, opts ...grpc.CallOption) (*GetMeResponse, error) {
	out := new(GetMeResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/GetMe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error) {
	out := new(RegisterUserResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/RegisterUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) PatchUser(ctx context.Context, in *PatchUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/PatchUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserByUsername(ctx context.Context, in *GetOneByUsernameMessage, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/GetUserByUsername", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetPublicUsers(ctx context.Context, in *GetMeMessage, opts ...grpc.CallOption) (*GetAllResponse, error) {
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/GetPublicUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CheckIfUserIsPrivate(ctx context.Context, in *GetOneMessage, opts ...grpc.CallOption) (*IsPrivateResponse, error) {
	out := new(IsPrivateResponse)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/CheckIfUserIsPrivate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetNotificationSettings(ctx context.Context, in *GetOneMessage, opts ...grpc.CallOption) (*NotificationSettings, error) {
	out := new(NotificationSettings)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/GetNotificationSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateNotificationSettings(ctx context.Context, in *UpdateNotificationSettingsRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/user_service_proto.UserService/UpdateNotificationSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	GetAll(context.Context, *SearchMessage) (*GetAllResponse, error)
	GetOne(context.Context, *GetOneMessage) (*UserResponse, error)
	GetMe(context.Context, *GetMeMessage) (*GetMeResponse, error)
	RegisterUser(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*UserResponse, error)
	PatchUser(context.Context, *PatchUserRequest) (*UserResponse, error)
	GetUserByUsername(context.Context, *GetOneByUsernameMessage) (*UserResponse, error)
	GetPublicUsers(context.Context, *GetMeMessage) (*GetAllResponse, error)
	CheckIfUserIsPrivate(context.Context, *GetOneMessage) (*IsPrivateResponse, error)
	GetNotificationSettings(context.Context, *GetOneMessage) (*NotificationSettings, error)
	UpdateNotificationSettings(context.Context, *UpdateNotificationSettingsRequest) (*EmptyMessage, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) GetAll(context.Context, *SearchMessage) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedUserServiceServer) GetOne(context.Context, *GetOneMessage) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedUserServiceServer) GetMe(context.Context, *GetMeMessage) (*GetMeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMe not implemented")
}
func (UnimplementedUserServiceServer) RegisterUser(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedUserServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserServiceServer) PatchUser(context.Context, *PatchUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PatchUser not implemented")
}
func (UnimplementedUserServiceServer) GetUserByUsername(context.Context, *GetOneByUsernameMessage) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByUsername not implemented")
}
func (UnimplementedUserServiceServer) GetPublicUsers(context.Context, *GetMeMessage) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPublicUsers not implemented")
}
func (UnimplementedUserServiceServer) CheckIfUserIsPrivate(context.Context, *GetOneMessage) (*IsPrivateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckIfUserIsPrivate not implemented")
}
func (UnimplementedUserServiceServer) GetNotificationSettings(context.Context, *GetOneMessage) (*NotificationSettings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotificationSettings not implemented")
}
func (UnimplementedUserServiceServer) UpdateNotificationSettings(context.Context, *UpdateNotificationSettingsRequest) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNotificationSettings not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetAll(ctx, req.(*SearchMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOneMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetOne(ctx, req.(*GetOneMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetMe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMeMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetMe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/GetMe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetMe(ctx, req.(*GetMeMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/RegisterUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RegisterUser(ctx, req.(*RegisterUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_PatchUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PatchUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).PatchUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/PatchUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).PatchUser(ctx, req.(*PatchUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserByUsername_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOneByUsernameMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserByUsername(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/GetUserByUsername",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserByUsername(ctx, req.(*GetOneByUsernameMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetPublicUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMeMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetPublicUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/GetPublicUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetPublicUsers(ctx, req.(*GetMeMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CheckIfUserIsPrivate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOneMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CheckIfUserIsPrivate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/CheckIfUserIsPrivate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CheckIfUserIsPrivate(ctx, req.(*GetOneMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetNotificationSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOneMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetNotificationSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/GetNotificationSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetNotificationSettings(ctx, req.(*GetOneMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateNotificationSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNotificationSettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateNotificationSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service_proto.UserService/UpdateNotificationSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateNotificationSettings(ctx, req.(*UpdateNotificationSettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user_service_proto.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAll",
			Handler:    _UserService_GetAll_Handler,
		},
		{
			MethodName: "GetOne",
			Handler:    _UserService_GetOne_Handler,
		},
		{
			MethodName: "GetMe",
			Handler:    _UserService_GetMe_Handler,
		},
		{
			MethodName: "RegisterUser",
			Handler:    _UserService_RegisterUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "PatchUser",
			Handler:    _UserService_PatchUser_Handler,
		},
		{
			MethodName: "GetUserByUsername",
			Handler:    _UserService_GetUserByUsername_Handler,
		},
		{
			MethodName: "GetPublicUsers",
			Handler:    _UserService_GetPublicUsers_Handler,
		},
		{
			MethodName: "CheckIfUserIsPrivate",
			Handler:    _UserService_CheckIfUserIsPrivate_Handler,
		},
		{
			MethodName: "GetNotificationSettings",
			Handler:    _UserService_GetNotificationSettings_Handler,
		},
		{
			MethodName: "UpdateNotificationSettings",
			Handler:    _UserService_UpdateNotificationSettings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user_service.proto",
}
