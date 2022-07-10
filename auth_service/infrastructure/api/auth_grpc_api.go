package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/tracer"

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
	span := tracer.StartSpanFromContext(ctx, "AuthenticateUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	loginRequest := mapLoginRequest(request.UserData)
	token, err := handler.service.AuthenticateUser(ctx, loginRequest)
	if err != nil {
		return nil, err
	}
	return &pb.JwtTokenResponse{
		Jwt: mapJwtToken(token),
	}, nil
}

func (handler *AuthHandler) AuthenticateTwoFactoryUser(ctx context.Context, request *pb.LoginTwoFactoryRequest) (*pb.JwtTokenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "AuthenticateTwoFactoryUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	token, err := handler.service.AuthenticateTwoFactoryUser(ctx, request)
	if err != nil {
		return nil, err
	}
	return &pb.JwtTokenResponse{
		Jwt: mapJwtToken(token),
	}, nil
}

func (handler *AuthHandler) GenerateTwoFactoryCode(ctx context.Context, request *pb.TwoFactoryLoginForCode) (*pb.TwoFactoryCode, error) {
	span := tracer.StartSpanFromContext(ctx, "GenerateTwoFactoryCodeAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	code, err := handler.service.GenerateTwoFactoryCode(ctx, request)
	if err != nil {
		return nil, err
	}
	return &pb.TwoFactoryCode{
		Code: code,
	}, nil
}

func (handler *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ValidateTokenAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	token := mapValidationRequest(req)
	claims, err := handler.service.ValidateToken(ctx, token)
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
	span := tracer.StartSpanFromContext(ctx, "PasswordLessLoginAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.PasswordlessLogin(ctx, req)
}

func (handler *AuthHandler) ConfirmEmailLogin(ctx context.Context, req *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ConfirmEmailAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.ConfirmEmailLogin(ctx, req)
}

func (handler *AuthHandler) ActivateAccount(ctx context.Context, req *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ActivateAccountAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.ActivateAccount(ctx, req)
}

func (handler *AuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ChangePasswordAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.ChangePassword(ctx, req)
}

func (handler *AuthHandler) RecoverAccount(ctx context.Context, req *pb.RecoverAccountRequest) (*pb.RecoverAccountResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "RecoverAccountAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.RecoverAccount(ctx, req)
}

func (handler *AuthHandler) SendAccountRecoveryMail(ctx context.Context, req *pb.AccountRecoveryMailRequest) (*pb.AccountRecoveryMailResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "SendAccountRecoveryMailAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.SendAccountRecoveryMail(ctx, req)
}

func (handler *AuthHandler) CreateNewAPIToken(ctx context.Context, request *pb.APITokenRequest) (*pb.NewAPITokenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateNewApiTokenAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("AuthHandler CreateNewAPIToken")
	fmt.Println(request.Username)
	return handler.service.GenerateAPIToken(ctx, request)
}

func (handler *AuthHandler) CheckApiToken(ctx context.Context, request *pb.JobPostingDtoRequest) (*pb.JobPostingDtoResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "CheckApiTokenAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.ValidateApiTokenFunc(ctx, request)
}
