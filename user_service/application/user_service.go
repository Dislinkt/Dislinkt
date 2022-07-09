package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"

	"github.com/dislinkt/user_service/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	store                    domain.UserStore
	registerUserOrchestrator *RegisterUserOrchestrator
	updateUserOrchestrator   *UpdateUserOrchestrator
	patchOrchestrator        *PatchUserOrchestrator
}

func NewUserService(store domain.UserStore, registerUserOrchestrator *RegisterUserOrchestrator,
	updateUserOrchestrator *UpdateUserOrchestrator,
	patchOrchestrator *PatchUserOrchestrator) *UserService {
	return &UserService{
		store:                    store,
		registerUserOrchestrator: registerUserOrchestrator,
		updateUserOrchestrator:   updateUserOrchestrator,
		patchOrchestrator:        patchOrchestrator,
	}
}

func (service *UserService) Register(user *domain.User) error {
	// span := tracer.StartSpanFromContext(ctx, "Register-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return errors.New("Invalid user data")
	}
	err := service.registerUserOrchestrator.Start(user)
	if err != nil {
		return err
	}

	return err
}
func (service *UserService) StartUpdate(user *domain.User) (*domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "Register-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return nil, errors.New("Invalid user data")
	}
	dbUser, err := service.store.FindByID(user.Id)
	if err != nil {
		return nil, err
	}

	err = service.updateUserOrchestrator.Start(user)
	if err != nil {
		return nil, err
	}

	return dbUser, err
}

func (service *UserService) Insert(user *domain.User) error {
	// span := tracer.StartSpanFromContext(ctx, "Register-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return errors.New("Invalid user data")
	}
	err := service.store.Insert(user)
	return err
}
func (service *UserService) Update(uuid uuid.UUID, user *domain.User) (*domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "Update-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	dbUser, err := service.store.FindByID(uuid)
	if err != nil {
		return nil, err
	}

	user.Id = uuid
	user.Email = dbUser.Email
	updatedUser, err := service.store.Update(user)
	if err != nil {
		return nil, err
	}
	return updatedUser, err
}

func (service *UserService) PatchUserStart(requestUser *domain.User) error {
	err := service.patchOrchestrator.Start(requestUser)
	if err != nil {
		return err
	}

	return err
}

func (service *UserService) PatchUser(updatePaths []string, requestUser *domain.User,
	username string) (*domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "Update-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	foundUser, err := service.store.FindByUsername(username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(foundUser)

	updatedUser, err := updateField(updatePaths, foundUser, requestUser)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dbUser, err := service.store.Update(updatedUser)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(dbUser)
	return dbUser, nil
}

func updateField(paths []string, user *domain.User, requestUser *domain.User) (*domain.User, error) {
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return nil, errors.New("Invalid user data")
	}

	if err := validator.New().Struct(requestUser); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return nil, errors.New("Invalid request user data")
	}

	for _, path := range paths {
		fmt.Println(path)
		switch path {
		case "private":
			user.Private = requestUser.Private
		case "name":
			user.Name = requestUser.Name
		case "surname":
			user.Surname = requestUser.Surname
		case "email":
			user.Email = requestUser.Email
		case "username":
			user.Username = requestUser.Username
		case "number":
			user.Number = requestUser.Number
		case "gender":
			user.Gender = requestUser.Gender
		case "date_of_birth":
			user.DateOfBirth = requestUser.DateOfBirth
		case "biography":
			user.Biography = requestUser.Biography
		default:
			return nil, status.Errorf(codes.PermissionDenied, "cannot update field '"+path+"'")
		}
	}
	user.UpdatedAt = time.Now()
	return user, nil
}

func (service *UserService) GetAll() (*[]domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "GetAll-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAll()
}

func (service *UserService) Search(searchText string) (*[]domain.User, error) {
	return service.store.Search(searchText)
}

func (service *UserService) GetOne(uuid uuid.UUID) (*domain.User, error) {
	return service.store.FindByID(uuid)
}

func (service *UserService) FindByUsername(username string) (*domain.User, error) {
	return service.store.FindByUsername(username)
}

func (service *UserService) Delete(user *domain.User) interface{} {
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return err
	}
	return service.store.Delete(user)
}

func (service *UserService) GetPublicUsers() (*[]domain.User, error) {
	return service.store.GetPublicUsers()
}

func (service *UserService) FindByID(uuid uuid.UUID) (*domain.User, error) {
	return service.store.FindByID(uuid)
}

func (service *UserService) UpdateNotificationSettings(uuid uuid.UUID, ConnectionNotifications bool, MessageNotifications bool, PostNotifications bool) error {
	return service.store.UpdateNotificationSettings(uuid, ConnectionNotifications, MessageNotifications, PostNotifications)
}
