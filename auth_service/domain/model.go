package domain

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id       uuid.UUID `gorm:"index:idx_name,unique"`
	Username string    `gorm:"unique"`
	Password string
	UserRole int
}

type LoginRequest struct {
	Username string `bson:"password"`
	Password string `bson:"password"`
}
