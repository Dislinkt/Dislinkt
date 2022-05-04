package application

import (
	"errors"
	"github.com/dislinkt/auth-service/domain"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) GetByUsername(username string) (*domain.User, error) {
	user, err := service.store.GetByUsername(username)

	if err != nil {
		return nil, errors.New("invalid user")
	}

	return user, err
}
