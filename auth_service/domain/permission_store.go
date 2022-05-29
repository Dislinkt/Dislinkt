package domain

type PermissionStore interface {
	Insert(permission *Permission) error
	GetAll() (*[]Permission, error)
	GetAllByRole(role int) (*[]Permission, error)
}
