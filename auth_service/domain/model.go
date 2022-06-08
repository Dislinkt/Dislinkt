package domain

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id        uuid.UUID `gorm:"index:idx_name,unique"`
	Username  string    `gorm:"unique" validate:"alphanum"`
	Email     string    `gorm:"unique" validate:"email"`
	Password  string
	UserRole  int
	Active    bool
	ApiToken  *string `gorm:"unique"`
	TotpToken string  `gorm:"totp_token"`
}

type LoginRequest struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type Permission struct {
	Id   uint `gorm:"primaryKey;auto_increment:true"`
	Role int
	Name string `gorm:"unique"`
}
