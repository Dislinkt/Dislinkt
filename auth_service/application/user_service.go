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

func (service *UserService) GetByEmail(email string) (*domain.User, error) {
	user, err := service.store.GetByEmail(email)

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

	// TODO: obrisala jer mi treba da imaju isti id da mogu da menjam username
	// newUUID := uuid.NewV4()
	// user.Id = newUUID
	err := service.store.Insert(user)
	return user.Id, err
}

func (service *UserService) Delete(user *domain.User) error {
	return service.store.Delete(user)
}

func (service *UserService) Update(uuid uuid.UUID, user *domain.User) error {
	// span := tracer.StartSpanFromContext(ctx, "Update-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	_, err := service.store.FindByID(uuid)
	if err != nil {
		return err
	}

	user.Id = uuid
	return service.store.Update(user)
}

func (service *UserService) GetById(uuid uuid.UUID) (*domain.User, error) {
	user, err := service.store.FindByID(uuid)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (service *UserService) ChangeUsername(stringId string, username string) error {

	id, err := uuid.FromString(stringId)
	user, err := service.GetById(id)
	if err != nil {
		return err
	}
	user.Username = username
	err = service.Update(id, user)
	if err != nil {
		return err
	}
	return err
}
