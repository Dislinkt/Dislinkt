package api

import (
	"fmt"
	"time"

	"github.com/dislinkt/common/saga/events"
	"github.com/dislinkt/user_service/domain"
	uuid "github.com/satori/go.uuid"
)

func mapCommandUser(command *events.RegisterUserCommand) *domain.User {
	fmt.Println("COMMAND: " + command.User.Password)
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
	fmt.Println("DOMAIN: " + userD.Password)
	return userD
}
