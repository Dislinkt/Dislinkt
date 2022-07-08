package domain

import (
	"time"
)

type UserNode struct {
	UserUID string
	Status  ProfileStatus
}

type ProfileStatus string

const (
	Private ProfileStatus = "PRIVATE"
	Public                = "PUBLIC"
)

type JobOffer struct {
	Id            string
	Position      string
	Preconditions string
	DatePosted    time.Time
	Duration      string
	Location      string
	Title         string
	Field         string
	Description   string
}
