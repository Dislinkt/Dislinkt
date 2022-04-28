package application

import (
	"fmt"
	"time"

	"github.com/dislinkt/user-service/domain"
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

func (service *UserService) UpdatePrivacy(privacy bool, uuid uuid.UUID) (*domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "Update-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	user, err := service.store.Find(uuid)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	user.Private = privacy
	user.UpdatedAt = time.Now()
	err = service.store.Update(user)
	return user, nil
}

func (service *UserService) GetAll() (*[]domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "GetAll-Service")
	// defer span.Finish()
	//
	// newCtx := tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAll()
}
