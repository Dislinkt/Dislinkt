package domain

import (
	"time"

	"github.com/gofrs/uuid"

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
	Name        string    ``
	Surname     string    ``
	Username    *string   `gorm:"unique" `
	Email       *string   `gorm:"unique" `
	Number      string
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
