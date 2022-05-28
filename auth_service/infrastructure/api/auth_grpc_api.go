package api

import (
	"context"
	//"github.com/dislinkt/auth_service/domain"
	"net/http"

	"github.com/dislinkt/auth_service/application"
	pb "github.com/dislinkt/common/proto/auth_service"
)

type AuthHandler struct {
	service *application.AuthService
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (handler *AuthHandler) AuthenticateUser(ctx context.Context, request *pb.LoginRequest) (*pb.JwtTokenResponse, error) {
	loginRequest := mapLoginRequest(request.UserData)
	token, err := handler.service.AuthenticateUser(loginRequest)
	if err != nil {
		return nil, err
	}
	return &pb.JwtTokenResponse{
		Jwt: mapJwtToken(token),
	}, nil
}

func (handler *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	token := mapValidationRequest(req)
	claims, err := handler.service.ValidateToken(token)
	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		}, nil
	}

	if claims != nil {
		return &pb.ValidateResponse{
			Status:   http.StatusOK,
			Username: username,
			Role:     role,
		}, nil
	}
	return &pb.ValidateResponse{
		Status:   http.StatusOK,
		Username: username,
		Role:     role,
	}, nil
}

func (handler *AuthHandler) PasswordlessLogin(ctx context.Context, req *pb.PasswordlessLoginRequest) (*pb.PasswordlessLoginResponse, error) {
	return handler.service.PasswordlessLogin(ctx, req)
}

func (handler *AuthHandler) ConfirmEmailLogin(ctx context.Context, req *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {
	return handler.service.ConfirmEmailLogin(ctx, req)
}

func (handler *AuthHandler) ActivateAccount(ctx context.Context, req *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	return handler.service.ActivateAccount(ctx, req)
}

func (handler *AuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return handler.service.ChangePassword(ctx, req)
}

func (handler *AuthHandler) RecoverAccount(ctx context.Context, req *pb.RecoverAccountRequest) (*pb.RecoverAccountResponse, error) {
	return handler.service.RecoverAccount(ctx, req)
}

func (handler *AuthHandler) SendAccountRecoveryMail(ctx context.Context, req *pb.AccountRecoveryMailRequest) (*pb.AccountRecoveryMailResponse, error) {
	return handler.service.SendAccountRecoveryMail(ctx, req)
}
