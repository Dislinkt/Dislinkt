package events

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeleteEducationCommandType int8

type EducationDelete struct {
	Id string `
	School       string             `
	Degree       string
	FieldOfStudy string
	StartDate    primitive.DateTime
	EndDate      primitive.DateTime
}

const (
	DeleteEducationInAdditional DeleteEducationCommandType = iota
	RollbackDeleteEducationInAdditional
	DeleteEducationInGraph
	CancelEducationDelete
	ApproveEducationDeleting
	UnknownDeleteEducationCommand
)

type DeleteEducationCommand struct {
	Education EducationDelete
	UserId    string
	Type      DeleteEducationCommandType
}

type DeleteEducationReplyType int8

const (
	AdditionalServiceEducationDeleted DeleteEducationReplyType = iota
	AdditionalServiceEducationNotDeleted
	AdditionalEducationDeleteRolledBack
	GraphDatabaseEducationDeleted
	GraphDatabaseEducationNotDeleted
	UnknownDeletedEducationReply
)

type DeleteEducationReply struct {
	Education EducationDelete
	UserId    string
	Type      DeleteEducationReplyType
}
