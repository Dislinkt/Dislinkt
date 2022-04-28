package domain

import uuid "github.com/satori/go.uuid"

type UserStore interface {
	Insert(user *User) error
	Update(user *User) error
	GetAll() (*[]User, error)
	Find(uuid uuid.UUID) (*User, error)
}
