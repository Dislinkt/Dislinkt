package domain

type UserStore interface {
	Insert(user *User) error
	Update(user *User) error
	GetAll() (*[]User, error)
}
