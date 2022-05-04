package api

import (
	"github.com/dislinkt/auth-service/domain"
	pb "github.com/dislinkt/common/proto/auth_service"
)

//func mapUser(userD *domain.User) *pb.User {
//	userPb := &pb.User{
//		Id:       userD.Id,
//		Username: userD.Username,
//	}
//	return userPb
//}

func mapLoginRequest(userData *pb.UserData) *domain.LoginRequest {
	loginReq := &domain.LoginRequest{
		Username: userData.Username,
		Password: userData.Password,
	}
	return loginReq
}

func mapJwtToken(jwt string) *pb.JwtToken {
	token := &pb.JwtToken{
		Jwt: jwt,
	}
	return token
}

func mapValidationRequest(req *pb.ValidateRequest) string {
	token := req.Jwt.Jwt
	return token
}
