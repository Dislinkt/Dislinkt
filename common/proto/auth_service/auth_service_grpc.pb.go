// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: auth_service.proto

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

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	AuthenticateUser(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*JwtTokenResponse, error)
	AuthenticateTwoFactoryUser(ctx context.Context, in *LoginTwoFactoryRequest, opts ...grpc.CallOption) (*JwtTokenResponse, error)
	GenerateTwoFactoryCode(ctx context.Context, in *TwoFactoryLoginForCode, opts ...grpc.CallOption) (*TwoFactoryCode, error)
	ValidateToken(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
	PasswordlessLogin(ctx context.Context, in *PasswordlessLoginRequest, opts ...grpc.CallOption) (*PasswordlessLoginResponse, error)
	ConfirmEmailLogin(ctx context.Context, in *ConfirmEmailLoginRequest, opts ...grpc.CallOption) (*ConfirmEmailLoginResponse, error)
	ActivateAccount(ctx context.Context, in *ActivationRequest, opts ...grpc.CallOption) (*ActivationResponse, error)
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error)
	RecoverAccount(ctx context.Context, in *RecoverAccountRequest, opts ...grpc.CallOption) (*RecoverAccountResponse, error)
	SendAccountRecoveryMail(ctx context.Context, in *AccountRecoveryMailRequest, opts ...grpc.CallOption) (*AccountRecoveryMailResponse, error)
	CreateNewAPIToken(ctx context.Context, in *APITokenRequest, opts ...grpc.CallOption) (*NewAPITokenResponse, error)
	CheckApiToken(ctx context.Context, in *JobPostingDtoRequest, opts ...grpc.CallOption) (*JobPostingDtoResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) AuthenticateUser(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*JwtTokenResponse, error) {
	out := new(JwtTokenResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/AuthenticateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) AuthenticateTwoFactoryUser(ctx context.Context, in *LoginTwoFactoryRequest, opts ...grpc.CallOption) (*JwtTokenResponse, error) {
	out := new(JwtTokenResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/AuthenticateTwoFactoryUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) GenerateTwoFactoryCode(ctx context.Context, in *TwoFactoryLoginForCode, opts ...grpc.CallOption) (*TwoFactoryCode, error) {
	out := new(TwoFactoryCode)
	err := c.cc.Invoke(ctx, "/proto.AuthService/GenerateTwoFactoryCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ValidateToken(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/ValidateToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) PasswordlessLogin(ctx context.Context, in *PasswordlessLoginRequest, opts ...grpc.CallOption) (*PasswordlessLoginResponse, error) {
	out := new(PasswordlessLoginResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/PasswordlessLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ConfirmEmailLogin(ctx context.Context, in *ConfirmEmailLoginRequest, opts ...grpc.CallOption) (*ConfirmEmailLoginResponse, error) {
	out := new(ConfirmEmailLoginResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/ConfirmEmailLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ActivateAccount(ctx context.Context, in *ActivationRequest, opts ...grpc.CallOption) (*ActivationResponse, error) {
	out := new(ActivationResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/ActivateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error) {
	out := new(ChangePasswordResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) RecoverAccount(ctx context.Context, in *RecoverAccountRequest, opts ...grpc.CallOption) (*RecoverAccountResponse, error) {
	out := new(RecoverAccountResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/RecoverAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) SendAccountRecoveryMail(ctx context.Context, in *AccountRecoveryMailRequest, opts ...grpc.CallOption) (*AccountRecoveryMailResponse, error) {
	out := new(AccountRecoveryMailResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/SendAccountRecoveryMail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) CreateNewAPIToken(ctx context.Context, in *APITokenRequest, opts ...grpc.CallOption) (*NewAPITokenResponse, error) {
	out := new(NewAPITokenResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/CreateNewAPIToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) CheckApiToken(ctx context.Context, in *JobPostingDtoRequest, opts ...grpc.CallOption) (*JobPostingDtoResponse, error) {
	out := new(JobPostingDtoResponse)
	err := c.cc.Invoke(ctx, "/proto.AuthService/CheckApiToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	AuthenticateUser(context.Context, *LoginRequest) (*JwtTokenResponse, error)
	AuthenticateTwoFactoryUser(context.Context, *LoginTwoFactoryRequest) (*JwtTokenResponse, error)
	GenerateTwoFactoryCode(context.Context, *TwoFactoryLoginForCode) (*TwoFactoryCode, error)
	ValidateToken(context.Context, *ValidateRequest) (*ValidateResponse, error)
	PasswordlessLogin(context.Context, *PasswordlessLoginRequest) (*PasswordlessLoginResponse, error)
	ConfirmEmailLogin(context.Context, *ConfirmEmailLoginRequest) (*ConfirmEmailLoginResponse, error)
	ActivateAccount(context.Context, *ActivationRequest) (*ActivationResponse, error)
	ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error)
	RecoverAccount(context.Context, *RecoverAccountRequest) (*RecoverAccountResponse, error)
	SendAccountRecoveryMail(context.Context, *AccountRecoveryMailRequest) (*AccountRecoveryMailResponse, error)
	CreateNewAPIToken(context.Context, *APITokenRequest) (*NewAPITokenResponse, error)
	CheckApiToken(context.Context, *JobPostingDtoRequest) (*JobPostingDtoResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) AuthenticateUser(context.Context, *LoginRequest) (*JwtTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateUser not implemented")
}
func (UnimplementedAuthServiceServer) AuthenticateTwoFactoryUser(context.Context, *LoginTwoFactoryRequest) (*JwtTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateTwoFactoryUser not implemented")
}
func (UnimplementedAuthServiceServer) GenerateTwoFactoryCode(context.Context, *TwoFactoryLoginForCode) (*TwoFactoryCode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateTwoFactoryCode not implemented")
}
func (UnimplementedAuthServiceServer) ValidateToken(context.Context, *ValidateRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (UnimplementedAuthServiceServer) PasswordlessLogin(context.Context, *PasswordlessLoginRequest) (*PasswordlessLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PasswordlessLogin not implemented")
}
func (UnimplementedAuthServiceServer) ConfirmEmailLogin(context.Context, *ConfirmEmailLoginRequest) (*ConfirmEmailLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmEmailLogin not implemented")
}
func (UnimplementedAuthServiceServer) ActivateAccount(context.Context, *ActivationRequest) (*ActivationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateAccount not implemented")
}
func (UnimplementedAuthServiceServer) ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedAuthServiceServer) RecoverAccount(context.Context, *RecoverAccountRequest) (*RecoverAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverAccount not implemented")
}
func (UnimplementedAuthServiceServer) SendAccountRecoveryMail(context.Context, *AccountRecoveryMailRequest) (*AccountRecoveryMailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendAccountRecoveryMail not implemented")
}
func (UnimplementedAuthServiceServer) CreateNewAPIToken(context.Context, *APITokenRequest) (*NewAPITokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewAPIToken not implemented")
}
func (UnimplementedAuthServiceServer) CheckApiToken(context.Context, *JobPostingDtoRequest) (*JobPostingDtoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckApiToken not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_AuthenticateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AuthenticateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/AuthenticateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AuthenticateUser(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_AuthenticateTwoFactoryUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginTwoFactoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AuthenticateTwoFactoryUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/AuthenticateTwoFactoryUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AuthenticateTwoFactoryUser(ctx, req.(*LoginTwoFactoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_GenerateTwoFactoryCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TwoFactoryLoginForCode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GenerateTwoFactoryCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/GenerateTwoFactoryCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GenerateTwoFactoryCode(ctx, req.(*TwoFactoryLoginForCode))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ValidateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ValidateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/ValidateToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ValidateToken(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_PasswordlessLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordlessLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).PasswordlessLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/PasswordlessLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).PasswordlessLogin(ctx, req.(*PasswordlessLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ConfirmEmailLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmEmailLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ConfirmEmailLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/ConfirmEmailLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ConfirmEmailLogin(ctx, req.(*ConfirmEmailLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ActivateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ActivateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/ActivateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ActivateAccount(ctx, req.(*ActivationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ChangePassword(ctx, req.(*ChangePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_RecoverAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecoverAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).RecoverAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/RecoverAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).RecoverAccount(ctx, req.(*RecoverAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_SendAccountRecoveryMail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountRecoveryMailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).SendAccountRecoveryMail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/SendAccountRecoveryMail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).SendAccountRecoveryMail(ctx, req.(*AccountRecoveryMailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_CreateNewAPIToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(APITokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).CreateNewAPIToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/CreateNewAPIToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).CreateNewAPIToken(ctx, req.(*APITokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_CheckApiToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobPostingDtoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).CheckApiToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuthService/CheckApiToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).CheckApiToken(ctx, req.(*JobPostingDtoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthenticateUser",
			Handler:    _AuthService_AuthenticateUser_Handler,
		},
		{
			MethodName: "AuthenticateTwoFactoryUser",
			Handler:    _AuthService_AuthenticateTwoFactoryUser_Handler,
		},
		{
			MethodName: "GenerateTwoFactoryCode",
			Handler:    _AuthService_GenerateTwoFactoryCode_Handler,
		},
		{
			MethodName: "ValidateToken",
			Handler:    _AuthService_ValidateToken_Handler,
		},
		{
			MethodName: "PasswordlessLogin",
			Handler:    _AuthService_PasswordlessLogin_Handler,
		},
		{
			MethodName: "ConfirmEmailLogin",
			Handler:    _AuthService_ConfirmEmailLogin_Handler,
		},
		{
			MethodName: "ActivateAccount",
			Handler:    _AuthService_ActivateAccount_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _AuthService_ChangePassword_Handler,
		},
		{
			MethodName: "RecoverAccount",
			Handler:    _AuthService_RecoverAccount_Handler,
		},
		{
			MethodName: "SendAccountRecoveryMail",
			Handler:    _AuthService_SendAccountRecoveryMail_Handler,
		},
		{
			MethodName: "CreateNewAPIToken",
			Handler:    _AuthService_CreateNewAPIToken_Handler,
		},
		{
			MethodName: "CheckApiToken",
			Handler:    _AuthService_CheckApiToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth_service.proto",
}
