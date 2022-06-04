package api

import (
	"log"

	"github.com/dislinkt/auth_service/domain"
	"github.com/dislinkt/common/saga/events"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func mapCommandUser(command *events.RegisterUserCommand) *domain.User {
	id, _ := uuid.FromString(command.User.Id)
	hashAndSalt, err := HashAndSaltPasswordIfStrongAndMatching(command.User.Password)
	if err != nil {
		return nil
	}
	// TODO: ACTIVE NA FALSE!!
	userD := &domain.User{
		Id:       id,
		Username: command.User.Username,
		Password: hashAndSalt,
		Email:    command.User.Email,
		UserRole: int(command.User.UserRole),
		Active:   true,
		ApiToken: nil,
	}
	return userD
}

func HashAndSaltPasswordIfStrongAndMatching(password string) (string, error) {
	// isWeak, _ := regexp.MatchString("^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[^!@#$%^&*(),.?\":{}|<>~'_+=]*)$", password)
	//
	// if isWeak {
	//	return "", errors.New("Password must contain minimum eight characters, at least one capital letter, one number and one special character")
	// }
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash), err
}
