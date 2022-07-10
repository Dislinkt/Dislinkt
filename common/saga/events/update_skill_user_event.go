package events

type UpdateSkillCommandType int8

type SkillUpdate struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

const (
	UpdateSkillInAdditional UpdateSkillCommandType = iota
	RollbackSkillUpdateInAdditional
	UpdateSkillInGraph
	CancelSkillUpdate
	ApproveSkillUpdating
	UnknownUpdateSkillCommand
)

type UpdateSkillCommand struct {
	Skill   SkillUpdate
	UserId  string
	Type    UpdateSkillCommandType
	OldName string
}

type UpdateSkillReplyType int8

const (
	AdditionalServiceSkillUpdated UpdateSkillReplyType = iota
	AdditionalServiceSkillNotUpdated
	AdditionalSkillUpdateRolledBack
	GraphDatabaseSkillUpdated
	GraphDatabaseSkillNotUpdated
	UnknownUpdatedSkillReply
)

type UpdateSkillReply struct {
	Skill   SkillUpdate
	UserId  string
	Type    UpdateSkillReplyType
	OldName string
}
