package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/dislinkt/common/tracer"
	"github.com/go-playground/validator/v10"
	"github.com/pquerna/otp/totp"
	"log"
	"net/smtp"
	"regexp"
	//	"github.com/nats-io/jwt/v2"
	"time"

	"github.com/dislinkt/auth_service/domain"
	"github.com/dislinkt/auth_service/startup/config"
	"github.com/dislinkt/common/interceptor"
	pb "github.com/dislinkt/common/proto/auth_service"
	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	Permissions []string `json:"permissions"`
}

func NewAuthService(userService *UserService, permissionStore domain.PermissionStore) *AuthService {
	return &AuthService{
		userService:     userService,
		permissionStore: permissionStore,
	}
}

func (auth *AuthService) AuthenticateUser(ctx context.Context, loginRequest *domain.LoginRequest) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "AuthenticateUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := auth.userService.GetByUsername(ctx, loginRequest.Username)
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
	token, err := auth.generateToken(ctx, user, expireTime)
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

func (auth *AuthService) generateToken(ctx context.Context, user *domain.User, expireTime int64) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "GenerateToken-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	//	rolesString, _ := json.Marshal(user.Roles)
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return "", errors.New("Invalid user data")
	}

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

func (auth *AuthService) AuthenticateTwoFactoryUser(ctx context.Context, loginRequest *pb.LoginTwoFactoryRequest) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "AuthenticateTwoFactoryUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := auth.userService.GetByUsername(ctx, loginRequest.Username)

	if !user.Active {
		return "", errors.New("user account not activated!")
	}

	valid := totp.Validate(loginRequest.Code, user.TotpToken)

	if !valid {
		return "", errors.New("Token not valid!")
	}

	expireTime := time.Now().Add(time.Hour).Unix()
	token, err := auth.generateToken(ctx, user, expireTime)
	if err != nil {
		return "", errors.New("invalid password")
	}

	return token, err
}

func (auth *AuthService) GenerateTwoFactoryCode(ctx context.Context, loginRequest *pb.TwoFactoryLoginForCode) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "GenerateTwoFactoryCode-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := auth.userService.GetByUsername(ctx, loginRequest.Username)

	if !user.Active {
		return "", errors.New("user account not activated!")
	}

	if !equalPasswords(user.Password, loginRequest.Password) {
		return "", errors.New("invalid password")
	}
	n := time.Now().UTC()
	code, err := totp.GenerateCode(user.TotpToken, n)

	if err != nil {
		return "", errors.New("Error generating token!")
	}

	return code, err
}

func (auth *AuthService) ValidateToken(ctx context.Context, signedToken string) (claims jwt.MapClaims, err error) {
	span := tracer.StartSpanFromContext(ctx, "ValidateToken-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
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
	span := tracer.StartSpanFromContext(ctx, "PasswordlessLogin-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	user, err := auth.userService.GetByEmail(ctx, request.Email)
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
	token, err := auth.generateToken(ctx, user, expireTime)
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
	urlRedirection := "https://localhost:4200/passwordless-login-validation/" + token

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
	span := tracer.StartSpanFromContext(ctx, "ConfirmEmailLogin-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	token, err := jwt.ParseWithClaims(request.Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
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
	span := tracer.StartSpanFromContext(ctx, "ChangePassword-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	user, err := auth.userService.GetByUsername(ctx, username)
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
	auth.userService.Update(ctx, user.Id, user)

	return &pb.ChangePasswordResponse{
		StatusCode: "200",
		Message:    "New password updated",
	}, nil
}

func HashAndSaltPasswordIfStrongAndMatching(password string) (string, error) {
	isStrong, _ := regexp.MatchString("[0-9A-Za-z!?#$@.*+_\\-]+", password)

	if !isStrong {
		return "", errors.New("Password not strong enough!")
	}
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash), err
}

func (auth *AuthService) SendActivationMail(ctx context.Context, username string) error {
	span := tracer.StartSpanFromContext(ctx, "SendActivationMail-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	user, err := auth.userService.GetByUsername(ctx, username)
	if err != nil || user == nil {
		return errors.New("invalid username")
	}

	expireTime := time.Now().Add(time.Hour).Unix()
	token, err := auth.generateToken(ctx, user, expireTime)

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
	urlRedirection := "https://localhost:4200/activate-account/" + token
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
	span := tracer.StartSpanFromContext(ctx, "ActivateAccount-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	token, err := jwt.ParseWithClaims(request.Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
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

	user, err := auth.userService.GetByUsername(ctx, claims.Username)
	if err != nil || user == nil {
		return nil, errors.New("invalid username")
	}

	user.Active = true
	auth.userService.Update(ctx, user.Id, user)

	return &pb.ActivationResponse{
		Token: request.Token,
	}, nil
}

func (auth *AuthService) SendAccountRecoveryMail(ctx context.Context, request *pb.AccountRecoveryMailRequest) (*pb.AccountRecoveryMailResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "SendAccountRecoveryMail-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	user, err := auth.userService.GetByEmail(ctx, request.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email")
	}

	expireTime := time.Now().Add(time.Hour).Unix()
	token, err := auth.generateToken(ctx, user, expireTime)

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
	urlRedirection := "https://localhost:4200/recover-account/" + token

	subject := "Subject: Account recovery\n"
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
	span := tracer.StartSpanFromContext(ctx, "RecoverAccount-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	token, err := jwt.ParseWithClaims(request.Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
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

	user, err := auth.userService.GetByUsername(ctx, claims.Username)
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
	auth.userService.Update(ctx, user.Id, user)

	return &pb.RecoverAccountResponse{
		StatusCode: "200",
		Message:    "User account recovered",
	}, nil
}

func (auth *AuthService) GenerateAPIToken(ctx context.Context, request *pb.APITokenRequest) (*pb.NewAPITokenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GenerateAPIToken-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("Auth Service GenerateAPIToken")
	var permissions []string
	permissions = append(permissions, "createJobOffer")
	user, err := auth.userService.GetByUsername(ctx, request.Username)
	fmt.Println(user)
	expireTime := time.Now().Add(time.Hour * 4).Unix()
	claims := ApiTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Id.String(),
			ExpiresAt: expireTime,
		},
		Username:    user.Username,
		Permissions: permissions,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	fmt.Println(config.NewConfig().PublicKey)
	jwtToken, err := token.SignedString([]byte(config.NewConfig().PublicKey))

	salted, err := HashAndSaltApiToken(jwtToken)
	user.ApiToken = &salted

	err = auth.userService.Update(ctx, user.Id, user)

	if err != nil {
		return nil, err
	}

	return &pb.NewAPITokenResponse{
		Token: jwtToken,
	}, nil
}

func (auth *AuthService) ValidateApiTokenFunc(ctx context.Context, request *pb.JobPostingDtoRequest) (*pb.JobPostingDtoResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ValidateApiTokenFunc-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	claims, err := auth.VerifyApiToken(ctx, request.ApiToken)
	if err != nil {
		return nil, nil
	}

	user, err := auth.userService.GetByUsername(ctx, claims.Username)
	if user == nil {
		fmt.Println("nema usera")
		return nil, nil
	}

	if !equalTokens(*user.ApiToken, request.ApiToken) {
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
		Token:         request.ApiToken,
	}, nil
}

func (auth *AuthService) VerifyApiToken(ctx context.Context, apiToken string) (claims *ApiTokenClaims, err error) {
	span := tracer.StartSpanFromContext(ctx, "VerifyApiToken-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	token, err := jwt.ParseWithClaims(apiToken, &ApiTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
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
