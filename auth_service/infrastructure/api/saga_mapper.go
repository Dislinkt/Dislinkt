package api

import (
	"github.com/dislinkt/auth_service/domain"
	events "github.com/dislinkt/common/saga/register_user"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func mapCommandUser(command *events.RegisterUserCommand) *domain.User {
	hashAndSalt, err := HashAndSaltPasswordIfStrongAndMatching(command.User.Password)
	if err != nil {
		return nil
	}
	userD := &domain.User{
		Id:       uuid.UUID{},
		Username: command.User.Username,
		Password: hashAndSalt,
		Email:    command.User.Email,
		UserRole: int(command.User.UserRole),
		Active:   false,
	}
	return userD
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
