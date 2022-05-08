package api

import (
	"github.com/dislinkt/auth_service/domain"
	events "github.com/dislinkt/common/saga/register_user"
	uuid "github.com/satori/go.uuid"
)

func mapCommandUser(command *events.RegisterUserCommand) *domain.User {
	userD := &domain.User{
		Id:       uuid.UUID{},
		Username: command.User.Username,
		Password: command.User.Password,
	}
	return userD
}
