package domain

import (
	uuid "github.com/gofrs/uuid"
	"time"

	"gorm.io/gorm"
)

type Role int

const (
	Regular Role = iota
	Admin
	Agent
)

type Gender int

const (
	Empty Gender = iota
	Male
	Female
)

type User struct {
	Id          uuid.UUID `gorm:"index:idx_name,unique"`
	Name        string    `validate:"alpha"`
	Surname     string    `validate:"alpha"`
	Username    *string   `gorm:"unique" validate:"alphanum"`
	Email       *string   `gorm:"unique" validate:"email"`
	Number      string    `validate:"numeric,omitempty"`
	Gender      Gender
	DateOfBirth string
	Password    string `gorm:"-"`
	UserRole    Role
	Biography   string
	Blocked     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Private     bool
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Username") {
		tx.Statement.SetColumn("Username", u.Username)
	}
	if tx.Statement.Changed() {
		tx.Statement.SetColumn("UpdatedAt", time.Now())
	}
	return nil
}
