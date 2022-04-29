package persistence

import (
	"github.com/dislinkt/user_service/domain"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserPostgresStore struct {
	db *gorm.DB
}

func NewUserPostgresStore(db *gorm.DB) (domain.UserStore, error) {
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}
	return &UserPostgresStore{
		db: db,
	}, nil
}

func (store *UserPostgresStore) Insert(user *domain.User) error {
	// span := tracer.StartSpanFromContext(ctx, "Insert-DB")
	// defer span.Finish()
	result := store.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *UserPostgresStore) Update(user *domain.User) error {
	// span := tracer.StartSpanFromContext(ctx, "Update-DB")
	// defer span.Finish()
	if result := store.db.Save(&user); result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *UserPostgresStore) GetAll() (*[]domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "GetAll-DB")
	// defer span.Finish()
	var users []domain.User
	result := store.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (store *UserPostgresStore) Find(uuid uuid.UUID) (user *domain.User, err error) {
	foundUser := domain.User{}
	if result := store.db.First(&foundUser, uuid); result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}
