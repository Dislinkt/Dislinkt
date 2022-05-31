package events

type UpdateUserCommandType int8

const (
	UpdateInUser UpdateUserCommandType = iota
	RollbackUpdateInUser
	UpdateInPost
	RollbackUpdateInPost
	UnknownUpdateCommand
	UserUpdateSucceeded
	UserUpdateCancelled
)

type UpdateUserCommand struct {
	User User
	Type UpdateUserCommandType
}

type UpdateUserReplyType int8

const (
	UserUpdatedInUser UpdateUserReplyType = iota
	UserNotUpdatedInUser
	UserRolledBackInUser
	UserUpdatedInPost
	UserNotUpdatedInPost
	UserRolledBackInPost
	UnknownUpdateReply
)

type UpdateUserReply struct {
	User User
	Type UpdateUserReplyType
}
