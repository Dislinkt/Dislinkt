package events

import "go.mongodb.org/mongo-driver/bson/primitive"

type Education struct {
	Id           primitive.ObjectID `bson:"_id"`
	School       string             `bson:"school"`
	Degree       string             `bson:"degree"`
	FieldOfStudy string             `bson:"field_of_study"`
	StartDate    primitive.DateTime `bson:"start_date"`
	EndDate      primitive.DateTime `bson:"end_date"`
}

type AddEducationCommandType int8

const (
	AddEducationInAdditional AddEducationCommandType = iota
	RollbackEducationInAdditional
	AddEducationInGraph
	CancelEducationAdd
	ApproveAdding
	UnknownAddEducationCommand
)

type AddEducationCommand struct {
	Education Education
	UserId    string
	Type      AddEducationCommandType
}

type AddEducationReplyType int8

const (
	AdditionalServiceAdded AddEducationReplyType = iota
	AdditionalServiceNotAdded
	AdditionalRolledBack
	GraphDatabaseAdded
	GraphDatabaseNotAdded
	UnknownAddEducationReply
)

type AddEducationReply struct {
	Education Education
	UserId    string
	Type      AddEducationReplyType
}
