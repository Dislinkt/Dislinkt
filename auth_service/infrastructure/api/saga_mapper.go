package api

import (
	"errors"
	"log"
	"regexp"

	"github.com/dislinkt/auth_service/domain"
	"github.com/dislinkt/common/saga/events"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func mapCommandUser(command *events.RegisterUserCommand) *domain.User {
	id, _ := uuid.FromString(command.User.Id)
	hashAndSalt, err := HashAndSaltPasswordIfStrongAndMatching(command.User.Password)
	if err != nil {
		return nil
	}
	userD := &domain.User{
		Id:        id,
		Username:  command.User.Username,
		Password:  hashAndSalt,
		Email:     command.User.Email,
		UserRole:  int(command.User.UserRole),
		Active:    false,
		ApiToken:  nil,
		TotpToken: nil,
		TotpQR:    nil,
	}
	return userD
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
