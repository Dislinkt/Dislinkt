package domain

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
	Id          string `gorm:"index:idx_name,unique"`
	Name        string
	Surname     string
	Username    string `gorm:"unique"`
	Email       string `gorm:"unique"`
	Number      string
	Gender      Gender
	DateOfBirth string
	Password    string
	UserRole    Role
	Biography   string
}
