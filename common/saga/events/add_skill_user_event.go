package events

import "go.mongodb.org/mongo-driver/bson/primitive"

type Skill struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type AddSkillCommandType int8

const (
	AddSkillInAdditional AddSkillCommandType = iota
	RollbackSkillInAdditional
	AddSkillInGraph
	CancelSkillAdd
	ApproveSkillAdding
	UnknownAddSkillCommand
)

type AddSkillCommand struct {
	Skill  Skill
	UserId string
	Type   AddSkillCommandType
}

type AddSkillReplyType int8

const (
	AdditionalServiceSkillAdded AddSkillReplyType = iota
	AdditionalServiceSkillNotAdded
	AdditionalSkillRolledBack
	GraphDatabaseSkillAdded
	GraphDatabaseSkillNotAdded
	UnknownAddSkillReply
)

type AddSkillReply struct {
	Skill  Skill
	UserId string
	Type   AddSkillReplyType
}
