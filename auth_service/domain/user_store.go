package domain

import uuid "github.com/gofrs/uuid"

type UserStore interface {
	Insert(user *User) error
	Update(user *User) error
	GetAll() (*[]User, error)
	FindByID(uuid uuid.UUID) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Delete(user *User) error
}
