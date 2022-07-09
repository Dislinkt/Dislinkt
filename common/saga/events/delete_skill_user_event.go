package events

type DeleteSkillCommandType int8

type SkillDelete struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

const (
	DeleteSkillInAdditional DeleteSkillCommandType = iota
	RollbackSkillDeleteInAdditional
	DeleteSkillInGraph
	CancelSkillDelete
	ApproveSkillDeleting
	UnknownDeleteSkillCommand
)

type DeleteSkillCommand struct {
	Skill  SkillDelete
	UserId string
	Type   DeleteSkillCommandType
}

type DeleteSkillReplyType int8

const (
	AdditionalServiceSkillDeleted DeleteSkillReplyType = iota
	AdditionalServiceSkillNotDeleted
	AdditionalSkillDeleteRolledBack
	GraphDatabaseSkillDeleted
	GraphDatabaseSkillNotDeleted
	UnknownDeletedSkillReply
)

type DeleteSkillReply struct {
	Skill  SkillDelete
	UserId string
	Type   DeleteSkillReplyType
}
