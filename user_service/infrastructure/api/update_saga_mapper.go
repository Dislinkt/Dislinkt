package api

import (
	"time"

	"github.com/dislinkt/common/saga/events"
	"github.com/dislinkt/user_service/domain"
	"github.com/gofrs/uuid"
)

func mapCommandUpdateUser(command *events.UpdateUserCommand) *domain.User {
	id, _ := uuid.FromString(command.User.Id)
	userD := &domain.User{
		Id:          id,
		Name:        command.User.Name,
		Surname:     command.User.Surname,
		Username:    &command.User.Username,
		Email:       &command.User.Email,
		Number:      command.User.Number,
		Gender:      domain.Gender(command.User.Gender),
		DateOfBirth: command.User.DateOfBirth,
		Biography:   command.User.Biography,
		UpdatedAt:   time.Now(),
		Private:     command.User.Private,
	}
	return userD
}
