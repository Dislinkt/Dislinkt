package api

import (
	"github.com/dislinkt/common/saga/events"
	"post_service/domain"
)

func mapPostCommandUpdateUser(command *events.UpdateUserCommand) *domain.User {
	userD := &domain.User{
		UserUUID: command.User.Id,
		Name:     command.User.Name,
		Surname:  command.User.Surname,
		Username: command.User.Username,
	}
	return userD
}
