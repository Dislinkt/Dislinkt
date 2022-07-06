package persistence

import (
	"encoding/json"

	"github.com/dislinkt/user_service/domain"
	"github.com/gofrs/uuid"
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
	// span := tracer.StartSpanFromContext(ctx, "Register-DB")
	// defer span.Finish()
	result := store.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *UserPostgresStore) Update(user *domain.User) (*domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "Update-DB")
	// defer span.Finish()
	var inInterface map[string]interface{}
	parsedUser, _ := json.Marshal(user)
	err1 := json.Unmarshal(parsedUser, &inInterface)
	if err1 != nil {
		return nil, err1
	}
	result := store.db.Model(&user).Updates(inInterface)
	if result.Error != nil {
		return nil, result.Error
	}
	dbUser, err := store.FindByID(user.Id)
	if err != nil {
		return nil, result.Error
	}
	return dbUser, nil
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

func (store *UserPostgresStore) GetPublicUsers() (*[]domain.User, error) {
	// span := tracer.StartSpanFromContext(ctx, "GetAll-DB")
	// defer span.Finish()
	var users []domain.User
	result := store.db.Where("private = false").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (store *UserPostgresStore) FindByID(uuid uuid.UUID) (user *domain.User, err error) {
	foundUser := domain.User{}
	if result := store.db.First(&foundUser, uuid); result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (store *UserPostgresStore) FindByUsername(username string) (user *domain.User, err error) {
	foundUser := domain.User{}
	if result := store.db.Where("username LIKE ?", username).Find(&foundUser); result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (store *UserPostgresStore) Search(searchText string) (*[]domain.User, error) {
	var users []domain.User
	arg := "%" + searchText + "%"
	result := store.db.Where("name LIKE ? OR surname LIKE ? OR username LIKE ? LIMIT 5", arg, arg, arg).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (store *UserPostgresStore) Delete(user *domain.User) error {
	result := store.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
