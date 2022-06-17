package api

import (
	"context"
	"fmt"

	// "github.com/dislinkt/auth_service/domain"
	"github.com/dislinkt/common/interceptor"
	logger "github.com/dislinkt/common/logging"

	//"github.com/dislinkt/auth_service/domain"
	"net/http"

	"github.com/dislinkt/auth_service/application"
	pb "github.com/dislinkt/common/proto/auth_service"
)

type AuthHandler struct {
	service *application.AuthService
	pb.UnimplementedAuthServiceServer
	logger *logger.Logger
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	logger := logger.InitLogger(context.TODO())
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

func (handler *AuthHandler) AuthenticateUser(ctx context.Context, request *pb.LoginRequest) (*pb.JwtTokenResponse, error) {
	handler.logger.InfoLogger.Infof("POST rr: AU {%s}", request.UserData.Username)
	loginRequest := mapLoginRequest(request.UserData)
	is2FA, token, err := handler.service.AuthenticateUser(loginRequest)
	if err != nil {
		handler.logger.WarnLogger.Warnf("BC {%s}", request.UserData.Username)
		return nil, err
	}
	return &pb.JwtTokenResponse{
		Jwt:   mapJwtToken(token),
		Is2FA: is2FA,
	}, nil
}

func (handler *AuthHandler) Get2FA(ctx context.Context, request *pb.Get2FARequest) (*pb.Get2FAResponse,
	error) {
	is2FA, err := handler.service.Get2FA(ctx, request)
	if err != nil {
		return nil, err
	}
	return &pb.Get2FAResponse{
		Is2FA: is2FA,
	}, nil
}

func (handler *AuthHandler) Set2FA(ctx context.Context, request *pb.Set2FARequest) (*pb.Set2FAResponse,
	error) {
	token, qr, err := handler.service.Set2FA(ctx, request)
	var returnToken string
	var returnQR string
	if token == nil {
		returnToken = ""
		returnQR = ""
	} else {
		returnToken = *token
		returnQR = convertByteToBase64(*qr)
	}
	if err != nil {
		return nil, err
	}
	return &pb.Set2FAResponse{
		TotpToken: returnToken,
		TotpQR:    returnQR,
	}, nil
}

func (handler *AuthHandler) AuthenticateTwoFactoryUser(ctx context.Context, request *pb.LoginTwoFactoryRequest) (*pb.JwtTokenResponse, error) {
	handler.logger.InfoLogger.Infof("POST rr: AU2F {%s}", request.Username)
	token, err := handler.service.AuthenticateTwoFactoryUser(request)
	if err != nil {
		handler.logger.WarnLogger.Warnf("BC {%s}", request.Username)
		return nil, err
	}
	return &pb.JwtTokenResponse{
		Jwt: mapJwtToken(token),
	}, nil
}

func (handler *AuthHandler) GenerateTwoFactoryCode(ctx context.Context, request *pb.TwoFactoryLoginForCode) (*pb.TwoFactoryCode, error) {
	handler.logger.InfoLogger.Infof("POST rr: GC {%s}", request.Username)
	code, err := handler.service.GenerateTwoFactoryCode(request)
	if err != nil {
		handler.logger.WarnLogger.Warnf("GC {%s}", request.Username)
		return nil, err
	}
	return &pb.TwoFactoryCode{
		Code: code,
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
	handler.logger.InfoLogger.Infof("PUT rr: PLG {%s}", req.Email)
	return handler.service.PasswordlessLogin(ctx, req)
}

func (handler *AuthHandler) ConfirmEmailLogin(ctx context.Context, req *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {
	return handler.service.ConfirmEmailLogin(ctx, req)
}

func (handler *AuthHandler) ActivateAccount(ctx context.Context, req *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	handler.logger.InfoLogger.Infof("PUT rr: AA {%s}", req.Token)
	return handler.service.ActivateAccount(ctx, req)
}

func (handler *AuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	handler.logger.InfoLogger.Infof("PUT rr: CP {%s}", username)
	return handler.service.ChangePassword(ctx, req)
}

func (handler *AuthHandler) RecoverAccount(ctx context.Context, req *pb.RecoverAccountRequest) (*pb.RecoverAccountResponse, error) {
	handler.logger.InfoLogger.Infof("POST rr: RA {%s}", req.Token)
	return handler.service.RecoverAccount(ctx, req)
}

func (handler *AuthHandler) SendAccountRecoveryMail(ctx context.Context, req *pb.AccountRecoveryMailRequest) (*pb.AccountRecoveryMailResponse, error) {
	return handler.service.SendAccountRecoveryMail(ctx, req)
}

func (handler *AuthHandler) CreateNewAPIToken(ctx context.Context, request *pb.APITokenRequest) (*pb.NewAPITokenResponse, error) {
	handler.logger.InfoLogger.Infof("GET rr: CAT {%s}", request.Username)
	fmt.Println("AuthHandler CreateNewAPIToken")
	fmt.Println(request.Username)
	return handler.service.GenerateAPIToken(ctx, request)
}

func (handler *AuthHandler) CheckApiToken(ctx context.Context, request *pb.JobPostingDtoRequest) (*pb.JobPostingDtoResponse, error) {
	return handler.service.ValidateApiTokenFunc(ctx, request)
}
