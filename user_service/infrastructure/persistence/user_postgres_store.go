package persistence

import (
	"context"
	"encoding/json"
	logger "github.com/dislinkt/common/logging"

	"github.com/dislinkt/user_service/domain"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserPostgresStore struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewUserPostgresStore(db *gorm.DB) (domain.UserStore, error) {
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}
	logger := logger.InitLogger(context.TODO())

	return &UserPostgresStore{
		db:     db,
		logger: logger,
	}, nil
}

func (store *UserPostgresStore) Insert(ctx context.Context, user *domain.User) error {
	result := store.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	store.logger.InfoLogger.Infof("UC {%s}", *user.Username)
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
	store.logger.InfoLogger.Infof("UU {%s}", *user.Username)
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

func (store *UserPostgresStore) FindByID(uuid uuid.UUID) (user *domain.User, err error) {
	foundUser := domain.User{}
	if result := store.db.First(&foundUser, uuid); result.Error != nil {
		store.logger.WarnLogger.Warnf("UNF {%s}", *user.Username)
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
	store.logger.InfoLogger.Infof("UD {%s}", *user.Username)
	return nil
}
