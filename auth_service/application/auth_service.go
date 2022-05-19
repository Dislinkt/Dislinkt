package application

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dislinkt/auth_service/domain"
	//	"github.com/nats-io/jwt/v2"
	"time"
)

type AuthService struct {
	userService *UserService
}

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

//type jwtClaims struct {
//	jwt.StandardClaims
//	Id    int64
//	Email string
//}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (auth *AuthService) AuthenticateUser(loginRequest *domain.LoginRequest) (string, error) {

	user, err := auth.userService.GetByUsername(loginRequest.Username)
	if err != nil || user == nil {
		return "", errors.New("invalid username")
	}

	if !equalPasswords(user.Password, loginRequest.Password) {
		return "", errors.New("invalid password")
	}

	expireTime := time.Now().Add(time.Hour).Unix() * 1000
	token, err := generateToken(user, expireTime)
	if err != nil {
		return "", errors.New("invalid password")
	}

	//rolesString, _ := json.Marshal(user.Roles)
	return token, err
}

func equalPasswords(password string, passwordRequest string) bool {

	//byteHash := []byte(hashedPwd)
	//plainPwd := []byte(passwordRequest)
	//err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	//if err != nil {
	//	return false
	//}

	if passwordRequest != password {
		return false
	}

	return true
}

func generateToken(user *domain.User, expireTime int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	//	rolesString, _ := json.Marshal(user.Roles)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["role"] = string(getRoleString(user.UserRole))
	claims["id"] = user.Id
	claims["exp"] = expireTime
	jwtToken, err := token.SignedString([]byte("Dislinkt"))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func getRoleString(role int) string {
	switch role {
	case 0:
		return "Regular"
	case 1:
		return "Admin"
	case 2:
		return "Agent"
	default:
		return "Regular"
	}
}

func (auth *AuthService) ValidateToken(signedToken string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//return []byte(os.Getenv("ACCESS_SECRET")), nil
		return []byte("Dislinkt"), nil
	})

	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("Couldn't parse claims")
	}

	//if !claims.VerifyExpiresAt(time.Now().Local().Unix()) {
	//	return nil, errors.New("JWT is expired")
	//}

	return claims, nil

}
