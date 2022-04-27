package application

import "github.com/dislinkt/user-service/domain"

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) Insert(user *domain.User) error {
	return service.Insert(user)
}

func (service *UserService) Update(user *domain.User) error {
	return service.store.Update(user)
}

func (service *UserService) GetAll() (*[]domain.User, error) {
	return service.store.GetAll()
}
