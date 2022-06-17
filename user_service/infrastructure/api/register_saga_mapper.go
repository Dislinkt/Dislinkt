package api

import (
	"time"

	"github.com/dislinkt/common/saga/events"
	"github.com/dislinkt/user_service/domain"
	"github.com/gofrs/uuid"
)

func mapCommandUser(command *events.RegisterUserCommand) *domain.User {
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
		UserRole:    0,
		Biography:   command.User.Biography,
		Blocked:     false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Private:     command.User.Private,
		Password:    command.User.Password,
	}
	return userD
}
