package domain

import (
	"google.golang.org/genproto/googleapis/type/date"
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
	DatePosted    date.Date
	Duration      string
	Location      string
	Title         string
	Field         string
}
