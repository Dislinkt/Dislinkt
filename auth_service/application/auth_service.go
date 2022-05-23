package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dislinkt/auth_service/domain"
	"github.com/dislinkt/auth_service/startup/config"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/smtp"

	pb "github.com/dislinkt/common/proto/auth_service"

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

	if !user.Active {
		return "", errors.New("user account not activated!")
	}

	if !equalPasswords(user.Password, loginRequest.Password) {
		return "", errors.New("invalid password")
	}

	expireTime := time.Now().Add(time.Hour).Unix()
	token, err := generateToken(user, expireTime)
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
	token, err := generateToken(user, expireTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not generate JWT token")
	}

	message := passwordlessLoginMailMessage(token)

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

func passwordlessLoginMailMessage(token string) []byte {
	urlRedirection := "http://localhost:3000/passwordless-login-validation/" + token

	subject := "Subject: Passwordless login\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body style=\"background-color: #f4f4f4; margin: 0 !important; padding: 0 !important;\">\n" +
		"    <!-- HIDDEN PREHEADER TEXT -->\n" +
		"    <div style=\"display: none; font-size: 1px; color: #fefefe; line-height: 1px; font-family: 'Lato', Helvetica, Arial, sans-serif; max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden;\"> We're thrilled to have you here! Get ready to dive into your new account.\n" +
		"    </div>\n" +
		"    <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\">\n" +
		"        <!-- LOGO -->\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#FFA73B\" align=\"center\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td align=\"center\" valign=\"top\" style=\"padding: 40px 10px 40px 10px;\"> </td>\n" +
		"                    </tr>\n" +
		"                </table>\n" +
		"            </td>\n" +
		"        </tr>\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#FFA73B\" align=\"center\" style=\"padding: 0px 10px 0px 10px;\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"center\" valign=\"top\" style=\"padding: 40px 20px 20px 20px; border-radius: 4px 4px 0px 0px; color: #111111; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 48px; font-weight: 400; letter-spacing: 4px; line-height: 48px;\">\n" +
		"                            <h1 style=\"font-size: 48px; font-weight: 400; margin: 2;\">Dislinkt</h1> <img src=\" https://img.icons8.com/cotton/100/000000/security-checked--v3.png\" width=\"125\" height=\"120\" style=\"display: block; border: 0px;\" />\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"                </table>\n" +
		"            </td>\n" +
		"        </tr>\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#f4f4f4\" align=\"center\" style=\"padding: 0px 10px 0px 10px;\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\" style=\"padding: 20px 30px 40px 30px; color: #666666; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 18px; font-weight: 400; line-height: 25px;\">\n" +
		"                            <p style=\"margin: 0;\">Someone tried to sign in to your account without password. Was that you?</p>\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\">\n" +
		"                            <table width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
		"                                <tr>\n" +
		"                                    <td bgcolor=\"#ffffff\" align=\"center\" style=\"padding: 20px 30px 60px 30px;\">\n" +
		"                                        <table border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
		"                                            <tr>\n" +
		"                                                <td align=\"center\" style=\"border-radius: 3px;\" bgcolor=\"#FFA73B\"><a href=\"" + urlRedirection + "\" target=\"_blank\" style=\"font-size: 20px; font-family: Helvetica, Arial, sans-serif; color: #ffffff; text-decoration: none; color: #ffffff; text-decoration: none; padding: 15px 25px; border-radius: 2px; border: 1px solid #FFA73B; display: inline-block;\">Yes! Login</a></td>\n" +
		"                                            </tr>\n" +
		"                                        </table>\n" +
		"                                    </td>\n" +
		"                                </tr>\n" +
		"                            </table>\n" +
		"                        </td>\n" +
		"                    </tr> \n" +
		"    </table>\n" +
		"    <br> <br>\n" +
		"</body>" +
		"</html>"
	message := []byte(subject + mime + body)
	return message
}

func (auth *AuthService) ConfirmEmailLogin(ctx context.Context, request *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {

	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("Dislinkt"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Couldn't parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("Couldn't parse claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Local().Unix(), true) {
		return nil, fmt.Errorf("JWT is expired")
	}

	if err != nil {
		return nil, fmt.Errorf("Invalid token: %w", err)
	}

	return &pb.ConfirmEmailLoginResponse{
		Token: request.Token,
	}, nil
}
