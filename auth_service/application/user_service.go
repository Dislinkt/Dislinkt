package application

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/dislinkt/common/validator"
	goValidator "github.com/go-playground/validator/v10"

	"github.com/dislinkt/auth_service/domain"
	"github.com/gofrs/uuid"
)

type UserService struct {
	store     domain.UserStore
	validator *goValidator.Validate
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store:     store,
		validator: validator.InitValidator(),
	}
}

func (service *UserService) GetByUsername(username string) (*domain.User, error) {
	if !isUsernameValid(username) {
		return nil, errors.New("unallowed characters in username")
	}
	fmt.Println(username)
	user, err := service.store.GetByUsername(username)

	if err != nil {
		return nil, errors.New("invalid user")
	}

	return user, err
}

func (service *UserService) GetByEmail(email string) (*domain.User, error) {
	if !isEmailValid(email) {
		return nil, errors.New("unallowed characters in email")
	}
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
	if err := service.validator.Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return uuid.Nil, errors.New("Invalid user data")
	}
	err := service.store.Insert(user)
	return user.Id, err
}

func (service *UserService) Delete(user *domain.User) error {
	if err := service.validator.Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return errors.New("Invalid user data")
	}
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
	if !isUsernameValid(username) {
		return errors.New("unallowed characters in username")
	}

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

func isUsernameValid(username string) bool {
	isValid, _ := regexp.MatchString("[0-9A-Za-z]+", username)

	return isValid
}

func isEmailValid(email string) bool {
	isValid, _ := regexp.MatchString("[a-z0-9.\\-_]{3,64}@([a-z0-9]+\\.){1,2}[a-z]+", email)

	return isValid
}
