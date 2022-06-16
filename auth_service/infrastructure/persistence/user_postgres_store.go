package persistence

import (
	"context"
	"fmt"
	"github.com/dislinkt/auth_service/domain"
	logger "github.com/dislinkt/common/logging"
	uuid "github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserPostgresStore struct {
	db     *gorm.DB
	logger *logger.Logger
}

func (store *UserPostgresStore) Delete(user *domain.User) error {
	if result := store.db.Delete(user); result.Error != nil {
		store.logger.WarnLogger.Warnf("UD {%s}", user.Username)
		return result.Error
	}
	return nil
}

func NewUserPostgresStore(db *gorm.DB) (domain.UserStore, error) {
	logger := logger.InitLogger(context.TODO())
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}
	return &UserPostgresStore{
		db:     db,
		logger: logger,
	}, nil
}

func (store *UserPostgresStore) Insert(user *domain.User) error {
	fmt.Println("TU SAM")
	result := store.db.Create(user)
	if result.Error != nil {
		store.logger.WarnLogger.Warnf("UC {%s}", user.Username)
		return result.Error
	}
	store.logger.InfoLogger.Infof("UC {%s}", user.Username)
	return nil
}

func (store *UserPostgresStore) Update(user *domain.User) error {
	if result := store.db.Save(&user); result.Error != nil {
		store.logger.WarnLogger.Warnf("UU {%s}", user.Username)
		return result.Error
	}
	store.logger.InfoLogger.Infof("UU {%s}", user.Username)
	return nil
}

func (store *UserPostgresStore) GetAll() (*[]domain.User, error) {
	var users []domain.User
	result := store.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (store *UserPostgresStore) GetByUsername(username string) (user *domain.User, err error) {
	foundUser := domain.User{}
	if result := store.db.Where("username=?", username).First(&foundUser); result.Error != nil {
		return nil, result.Error
	}
	fmt.Println(foundUser)
	fmt.Println(&foundUser)
	return &foundUser, nil
}

func (store *UserPostgresStore) FindByID(uuid uuid.UUID) (user *domain.User, err error) {
	foundUser := domain.User{}
	if result := store.db.First(&foundUser, uuid); result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (store *UserPostgresStore) GetByEmail(email string) (user *domain.User, err error) {
	foundUser := domain.User{}
	if result := store.db.Where("email=?", email).First(&foundUser); result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (store *UserPostgresStore) UpdateAPIToken(id, apiToken string) error {
	var auth domain.User
	err := store.db.First(&auth, "id = ?", id)
	store.db.Model(&domain.User{}).Where("Id = ?", id).Update("ApiToken", apiToken)
	if err != nil {
		return err.Error
	}
	return nil
}
