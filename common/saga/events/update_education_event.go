package events

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdateEducationCommandType int8

type EducationUpdate struct {
	Id           string
	School       string
	Degree       string
	FieldOfStudy string
	StartDate    primitive.DateTime
	EndDate      primitive.DateTime
}

const (
	UpdateEducationInAdditional UpdateEducationCommandType = iota
	RollbackEducationUpdateInAdditional
	UpdateEducationInGraph
	CancelEducationUpdate
	ApproveEducationUpdating
	UnknownUpdateEducationCommand
)

type UpdateEducationCommand struct {
	Education    EducationUpdate
	UserId       string
	Type         UpdateEducationCommandType
	OldFieldName string
}

type UpdateEducationReplyType int8

const (
	AdditionalServiceEducationUpdated UpdateEducationReplyType = iota
	AdditionalServiceEducationNotUpdated
	AdditionalEducationUpdateRolledBack
	GraphDatabaseEducationUpdated
	GraphDatabaseEducationNotUpdated
	UnknownUpdatedEducationReply
)

type UpdateEducationReply struct {
	Education    EducationUpdate
	UserId       string
	Type         UpdateEducationReplyType
	OldFieldName string
}
