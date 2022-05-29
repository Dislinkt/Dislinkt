package domain

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id       uuid.UUID `gorm:"index:idx_name,unique"`
	Username string    `gorm:"unique"`
	Email    string    `gorm:"unique"`
	Password string
	UserRole int
	Active   bool
}

type LoginRequest struct {
	Username string `bson:"password"`
	Password string `bson:"password"`
}

type Permission struct {
	Id   uint `gorm:"primaryKey;auto_increment:true"`
	Role int
	Name string `gorm:"unique"`
}
