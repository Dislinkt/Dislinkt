package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/dislinkt/common/tracer"
	"github.com/go-playground/validator/v10"
	"regexp"

	"github.com/dislinkt/auth_service/domain"
	uuid "github.com/gofrs/uuid"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetByUsername-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetByEmail-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !isEmailValid(email) {
		return nil, errors.New("unallowed characters in email")
	}
	user, err := service.store.GetByEmail(email)

	if err != nil {
		return nil, errors.New("invalid user")
	}

	return user, err
}

func (service *UserService) Insert(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	span := tracer.StartSpanFromContext(ctx, "Insert-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	// TODO: obrisala jer mi treba da imaju isti id da mogu da menjam username
	// newUUID := uuid.NewV4()
	// user.Id = newUUID
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return uuid.Nil, errors.New("Invalid user data")
	}
	err := service.store.Insert(user)
	return user.Id, err
}

func (service *UserService) Delete(ctx context.Context, user *domain.User) error {
	span := tracer.StartSpanFromContext(ctx, "Delete-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return errors.New("Invalid user data")
	}
	return service.store.Delete(user)
}

func (service *UserService) Update(ctx context.Context, uuid uuid.UUID, user *domain.User) error {
	span := tracer.StartSpanFromContext(ctx, "Update-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	_, err := service.store.FindByID(uuid)
	if err != nil {
		return err
	}

	user.Id = uuid
	return service.store.Update(user)
}

func (service *UserService) GetById(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetById-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := service.store.FindByID(uuid)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (service *UserService) ChangeUsername(ctx context.Context, stringId string, username string) error {
	span := tracer.StartSpanFromContext(ctx, "ChangeUsername-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if !isUsernameValid(username) {
		return errors.New("unallowed characters in username")
	}

	id, err := uuid.FromString(stringId)
	user, err := service.GetById(ctx, id)
	if err != nil {
		return err
	}
	user.Username = username
	err = service.Update(ctx, id, user)
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
