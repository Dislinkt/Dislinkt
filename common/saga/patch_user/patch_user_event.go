package patch_user

type User struct {
	Id       string
	Username string
	Private  bool
}

type PatchUserCommandType int8

const (
	PatchUserInUser PatchUserCommandType = iota
	PatchUserInConnection
	CancelPatch
	RollbackPatchInUser
	PatchSucceeded
	UnknownPatchCommand
)

type PatchUserCommand struct {
	User User
	Type PatchUserCommandType
}

type PatchUserReplyType int8

const (
	PatchedUserInUser PatchUserReplyType = iota
	PatchFailedInUser
	PatchFailedInConnection
	PatchedInConnection
	UnknownPatchReply
)

type PatchUserReply struct {
	User User
	Type PatchUserReplyType
}
