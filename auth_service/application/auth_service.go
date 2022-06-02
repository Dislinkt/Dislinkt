package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dislinkt/auth_service/domain"
	"github.com/dislinkt/auth_service/startup/config"
	"github.com/dislinkt/common/interceptor"
	pb "github.com/dislinkt/common/proto/auth_service"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/smtp"

	//	"github.com/nats-io/jwt/v2"
	"time"
)

type AuthService struct {
	userService     *UserService
	permissionStore domain.PermissionStore
}

type Claims struct {
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

type ApiTokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewAuthService(userService *UserService, permissionStore domain.PermissionStore) *AuthService {
	return &AuthService{
		userService:     userService,
		permissionStore: permissionStore,
	}
}

func (auth *AuthService) AuthenticateUser(loginRequest *domain.LoginRequest) (string, error) {

	user, err := auth.userService.GetByUsername(loginRequest.Username)
	if err != nil || user == nil {
		return "", errors.New("invalid username")
	}

	if !user.Active {
		return "", errors.New("user account not activated!")
	}

	if !equalPasswords(user.Password, loginRequest.Password) {
		return "", errors.New("invalid password")
	}

	expireTime := time.Now().Add(time.Hour).Unix()
	token, err := auth.generateToken(user, expireTime)
	if err != nil {
		return "", errors.New("invalid password")
	}

	return token, err
}

func equalPasswords(hashedPwd string, passwordRequest string) bool {

	byteHash := []byte(hashedPwd)
	plainPwd := []byte(passwordRequest)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (auth *AuthService) generateToken(user *domain.User, expireTime int64) (string, error) {
	//	rolesString, _ := json.Marshal(user.Roles)

	var permissionNames []string
	permissions, err := auth.permissionStore.GetAllByRole(user.UserRole)
	if err != nil {
		fmt.Println(err)
	}
	for _, permission := range *permissions {
		permissionNames = append(permissionNames, permission.Name)
	}

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Id.String(),
			ExpiresAt: expireTime,
		},
		Username:    user.Username,
		Role:        getRoleString(user.UserRole),
		Permissions: permissionNames,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	fmt.Println(config.NewConfig().PublicKey)
	jwtToken, err := token.SignedString([]byte(config.NewConfig().PublicKey))
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
		return []byte(config.NewConfig().PublicKey), nil
	})

	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("Couldn't parse claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Local().Unix(), true) {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil

}

func (auth *AuthService) PasswordlessLogin(ctx context.Context, request *pb.PasswordlessLoginRequest) (*pb.PasswordlessLoginResponse, error) {

	user, err := auth.userService.GetByEmail(request.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid username")
	}

	if !user.Active {
		return nil, errors.New("user account not activated!")
	}

	if err != nil {
		return nil, errors.New("there is no user with that email or account is not activated")
	}

	from := config.NewConfig().EmailSender
	password := config.NewConfig().EmailPassword

	to := []string{
		request.Email,
	}

	smtpHost := config.NewConfig().EmailHost
	smtpPort := config.NewConfig().EmailPort

	expireTime := time.Now().Add(time.Hour).Unix()
	token, err := auth.generateToken(user, expireTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not generate JWT token")
	}

	message := passwordlessLoginMailMessage(token, user.Username)

	// Authentication.
	authentication := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, authentication, from, to, message)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error while sending mail")
	}

	return &pb.PasswordlessLoginResponse{
		Success: "Email Sent Successfully! Check your email.",
	}, nil
}

func passwordlessLoginMailMessage(token string, username string) []byte {
	urlRedirection := "http://localhost:4200/passwordless-login-validation/" + token

	subject := "Subject: Passwordless login\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body>\n" +
		"Hello " + username + "! Click on link to log in: " + urlRedirection +
		"<br> <br>\n" +
		"</body>" +
		"</html>"
	message := []byte(subject + mime + body)
	return message
}

func (auth *AuthService) ConfirmEmailLogin(ctx context.Context, request *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {

	token, err := jwt.ParseWithClaims(request.Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.NewConfig().PublicKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, fmt.Errorf("Couldn't parse claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Local().Unix(), true) {
		return nil, fmt.Errorf("JWT is expired")
	}

	return &pb.ConfirmEmailLoginResponse{
		Token: request.Token,
	}, nil
}

func (auth *AuthService) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	username := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	user, err := auth.userService.GetByUsername(username)
	if err != nil {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "User not found",
		}, errors.New("User not found")
	}

	if request.NewPassword != request.NewReenteredPassword {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "New passwords do not match",
		}, errors.New("New passwords do not match")
	}

	oldMatched := equalPasswords(user.Password, request.OldPassword)
	if !oldMatched {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "Old password does not match",
		}, errors.New("Old password does not match")
	}

	hashedNewPassword, err := HashAndSaltPasswordIfStrongAndMatching(request.NewPassword)
	if err != nil || hashedNewPassword == "" {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	user.Password = hashedNewPassword
	auth.userService.Update(user.Id, user)

	return &pb.ChangePasswordResponse{
		StatusCode: "200",
		Message:    "New password updated",
	}, nil
}

func HashAndSaltPasswordIfStrongAndMatching(password string) (string, error) {
	//isWeak, _ := regexp.MatchString("^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[^!@#$%^&*(),.?\":{}|<>~'_+=]*)$", password)
	//
	//if isWeak {
	//	return "", errors.New("Password must contain minimum eight characters, at least one capital letter, one number and one special character")
	//}
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash), err
}

func (auth *AuthService) SendActivationMail(username string) error {
	user, err := auth.userService.GetByUsername(username)
	if err != nil || user == nil {
		return errors.New("invalid username")
	}

	expireTime := time.Now().Add(time.Hour).Unix()
	token, err := auth.generateToken(user, expireTime)

	message := verificationMailMessage(token, username)

	from := config.NewConfig().EmailSender
	emailPassword := config.NewConfig().EmailPassword
	to := []string{user.Email}

	host := config.NewConfig().EmailHost
	port := config.NewConfig().EmailPort
	smtpAddress := host + ":" + port
	authMail := smtp.PlainAuth("", from, emailPassword, host)
	errSendingMail := smtp.SendMail(smtpAddress, authMail, from, to, message)
	if errSendingMail != nil {
		fmt.Println("err:  ", errSendingMail)
		return errSendingMail
	}
	return nil
}

func verificationMailMessage(token string, username string) []byte {
	// TODO SD: port se moze izvuci iz env var - 4200
	urlRedirection := "http://localhost:4200/activate-account/" + token
	fmt.Println("MAIL MESSAGE")

	subject := "Subject: Account activation\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body>\n" +
		"Hello " + username + "! Please confirm your account with click on link: " + urlRedirection +
		"<br> <br>\n" +
		"</body>" +
		"</html>"
	message := []byte(subject + mime + body)
	return message
}

func (auth *AuthService) ActivateAccount(ctx context.Context, request *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	token, err := jwt.ParseWithClaims(request.Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.NewConfig().PublicKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, fmt.Errorf("Couldn't parse claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Local().Unix(), true) {
		return nil, fmt.Errorf("JWT is expired")
	}

	user, err := auth.userService.GetByUsername(claims.Username)
	if err != nil || user == nil {
		return nil, errors.New("invalid username")
	}

	user.Active = true
	auth.userService.Update(user.Id, user)

	return &pb.ActivationResponse{
		Token: request.Token,
	}, nil
}

func (auth *AuthService) SendAccountRecoveryMail(ctx context.Context, request *pb.AccountRecoveryMailRequest) (*pb.AccountRecoveryMailResponse, error) {
	user, err := auth.userService.GetByEmail(request.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email")
	}

	expireTime := time.Now().Add(time.Hour).Unix()
	token, err := auth.generateToken(user, expireTime)

	message := recoverAccountMailMessage(token, user.Username)

	from := config.NewConfig().EmailSender
	emailPassword := config.NewConfig().EmailPassword
	to := []string{user.Email}

	host := config.NewConfig().EmailHost
	port := config.NewConfig().EmailPort
	smtpAddress := host + ":" + port
	authMail := smtp.PlainAuth("", from, emailPassword, host)
	errSendingMail := smtp.SendMail(smtpAddress, authMail, from, to, message)
	if errSendingMail != nil {
		fmt.Println("err:  ", errSendingMail)
		return nil, errSendingMail
	}
	return &pb.AccountRecoveryMailResponse{
		Success: "Email Sent Successfully! Check your email.",
	}, nil
}

func recoverAccountMailMessage(token string, username string) []byte {
	// TODO SD: port se moze izvuci iz env var - 4200
	urlRedirection := "http://localhost:4200/recover-account/" + token

	subject := "Subject: Account activation\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body>\n" +
		"Hello " + username + "! Recover your account with click on link: " + urlRedirection +
		"<br> <br>\n" +
		"</body>" +
		"</html>"
	message := []byte(subject + mime + body)
	return message
}

func (auth *AuthService) RecoverAccount(ctx context.Context, request *pb.RecoverAccountRequest) (*pb.RecoverAccountResponse, error) {
	token, err := jwt.ParseWithClaims(request.Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.NewConfig().PublicKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, fmt.Errorf("Couldn't parse claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Local().Unix(), true) {
		return nil, fmt.Errorf("JWT is expired")
	}

	user, err := auth.userService.GetByUsername(claims.Username)
	if err != nil {
		return &pb.RecoverAccountResponse{
			StatusCode: "500",
			Message:    "User not found",
		}, errors.New("User not found")
	}

	if request.NewPassword != request.NewReenteredPassword {
		return &pb.RecoverAccountResponse{
			StatusCode: "500",
			Message:    "New passwords do not match",
		}, errors.New("New passwords do not match")
	}

	hashedNewPassword, err := HashAndSaltPasswordIfStrongAndMatching(request.NewPassword)
	if err != nil || hashedNewPassword == "" {
		return &pb.RecoverAccountResponse{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	user.Password = hashedNewPassword
	auth.userService.Update(user.Id, user)

	return &pb.RecoverAccountResponse{
		StatusCode: "200",
		Message:    "User account recovered",
	}, nil
}

func (auth *AuthService) GenerateAPIToken(ctx context.Context, request *pb.APITokenRequest) (*pb.NewAPITokenResponse, error) {
	fmt.Println("Auth Service GenerateAPIToken")
	user, err := auth.userService.GetByUsername(request.Username)
	fmt.Println(user)
	expireTime := time.Now().Add(time.Hour * 4).Unix()
	claims := ApiTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Id.String(),
			ExpiresAt: expireTime,
		},
		Username: user.Username,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	fmt.Println(config.NewConfig().PublicKey)
	jwtToken, err := token.SignedString([]byte(config.NewConfig().PublicKey))

	user.ApiToken, err = HashAndSaltApiToken(jwtToken)

	err = auth.userService.Update(user.Id, user)

	if err != nil {
		return nil, err
	}

	return &pb.NewAPITokenResponse{
		Token: jwtToken,
	}, nil
}

func (auth *AuthService) ValidateApiTokenFunc(ctx context.Context, request *pb.JobPostingDtoRequest) (*pb.JobPostingDtoResponse, error) {
	claims, err := auth.VerifyApiToken(request.ApiToken)
	if err != nil {
		return nil, nil
	}

	user, err := auth.userService.GetByUsername(claims.Username)
	if user == nil {
		fmt.Println("nema usera")
		return nil, nil
	}

	if !equalTokens(user.ApiToken, request.ApiToken) {
		fmt.Println("Greska hash")
		return nil, nil
	}

	return &pb.JobPostingDtoResponse{
		Position: &pb.EmployeePositionDto{
			Name:      request.Position.Name,
			Seniority: request.Position.Seniority,
		},
		Username:      user.Username,
		Message:       "Token found",
		Duration:      request.Duration,
		DatePosted:    request.DatePosted,
		Preconditions: request.Preconditions,
		Description:   request.Description,
	}, nil
}

func (interceptor *AuthService) VerifyApiToken(apiToken string) (claims *ApiTokenClaims, err error) {
	token, err := jwt.ParseWithClaims(apiToken, &ApiTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.NewConfig().PublicKey), nil
	})

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*ApiTokenClaims)

	if !ok {
		return nil, fmt.Errorf("Couldn't parse claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Local().Unix(), true) {
		return nil, fmt.Errorf("JWT is expired")
	}

	return claims, nil
}

//func (auth *AuthService) ValidateApiTokenFunc(ctx context.Context, request *pb.JobPostingDtoRequest) (*pb.JobPostingDtoResponse, error) {
//
//}

func HashAndSaltApiToken(apiToken string) (string, error) {

	pwd := []byte(apiToken)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash), err
}

func equalTokens(hashedTok string, tokenRequest string) bool {

	byteHash := []byte(hashedTok)
	plainPwd := []byte(tokenRequest)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
