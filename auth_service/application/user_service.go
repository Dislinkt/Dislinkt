package application

import (
	"errors"
	"fmt"
	"github.com/dislinkt/auth_service/domain"
	uuid "github.com/satori/go.uuid"
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
	fmt.Println(username)
	user, err := service.store.GetByUsername(username)

	if err != nil {
		return nil, errors.New("invalid user")
	}

	return user, err
}

func (service *UserService) Insert(user *domain.User) (uuid.UUID, error) {
	// span := tracer.StartSpanFromContext(ctx, "Register-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)

	newUUID := uuid.NewV4()
	user.Id = newUUID
	err := service.store.Insert(user)
	return newUUID, err
}

func (service *UserService) Delete(user *domain.User) error {
	return service.store.Delete(user)
}
