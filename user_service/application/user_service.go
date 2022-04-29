package application

import (
	"fmt"
	"time"

	"github.com/dislinkt/user_service/domain"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) Insert(user *domain.User) error {
	// span := tracer.StartSpanFromContext(ctx, "Insert-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)

	newUUID := uuid.NewV4()
	user.Id = newUUID
	err := service.store.Insert(user)

	if err != nil {
		fmt.Println("///////////////////////////")
		fmt.Println(err.Error())
		fmt.Println("///////////////////////////")
	}
	return err
}

func (service *UserService) Update(uuid uuid.UUID, user *domain.User) error {
	// span := tracer.StartSpanFromContext(ctx, "Update-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	_, err := service.store.Find(uuid)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	user.Id = uuid
	user.UpdatedAt = time.Now()
	return service.store.Update(user)
}

func (service *UserService) PatchUser(updatePaths []string, requestUser *domain.User,
	uuid uuid.UUID) (*domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "Update-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	foundUser, err := service.store.Find(uuid)
	fmt.Println(uuid.String())
	fmt.Println(foundUser.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	updatedUser, err := updateField(updatePaths, foundUser, requestUser)
	if err != nil {
		return nil, err
	}
	err = service.store.Update(updatedUser)
	return updatedUser, nil
}

func updateField(paths []string, user *domain.User, requestUser *domain.User) (*domain.User, error) {
	for _, path := range paths {
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
