package persistence

import (
	"fmt"
	"github.com/dislinkt/auth_service/domain"
	"gorm.io/gorm"
)

type PermissionPostgresStore struct {
	db *gorm.DB
}

func NewPermissionPostgresStore(db *gorm.DB) (domain.PermissionStore, error) {
	err := db.AutoMigrate(&domain.Permission{})
	if err != nil {
		return nil, err
	}
	return &PermissionPostgresStore{
		db: db,
	}, nil
}

func (store *PermissionPostgresStore) Insert(permission *domain.Permission) error {
	fmt.Println("TU SAM")
	result := store.db.Create(permission)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *PermissionPostgresStore) GetAll() (*[]domain.Permission, error) {
	var permissions []domain.Permission
	result := store.db.Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return &permissions, nil
}

func (store *PermissionPostgresStore) GetAllByRole(role int) (*[]domain.Permission, error) {
	var permissions []domain.Permission
	if result := store.db.Where("role = ?", role).Find(&permissions); result.Error != nil {
		return nil, result.Error
	}
	return &permissions, nil
}
