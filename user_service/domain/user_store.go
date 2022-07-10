package domain

import (
	"github.com/gofrs/uuid"
)

type UserStore interface {
	Insert(user *User) error
	Update(user *User) (*User, error)
	GetAll() (*[]User, error)
	FindByID(uuid uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
	Search(searchText string) (*[]User, error)
	Delete(user *User) error
	GetPublicUsers() (*[]User, error)
	UpdateNotificationSettings(uuid uuid.UUID, ConnectionNotifications bool, MessageNotifications bool, PostNotifications bool) error
}
