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

	//var user domain.User

	//if result := handler.userDB.Where(&domain.User{Email: claims.Email}).First(&user); result.Error != nil {
	//	return &pb.ValidateResponse{
	//		Status: http.StatusNotFound,
	//		Error:  "User not found",
	//	}, nil
	//}

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
