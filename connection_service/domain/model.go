package domain

type UserNode struct {
	UserUID string
	Status  ProfileStatus
}

type ProfileStatus string

const (
	Private ProfileStatus = "PRIVATE"
	Public                = "PUBLIC"
)
