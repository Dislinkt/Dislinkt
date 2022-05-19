package register_user

type Role int

const (
	Regular Role = iota
	Admin
	Agent
)

type Gender int

const (
	Empty Gender = iota
	Male
	Female
)

type User struct {
	Id          string
	Name        string
	Surname     string
	Username    string
	Email       string
	Number      string
	Gender      Gender
	DateOfBirth string
	Password    string
	UserRole    Role
	Biography   string
	Private     bool
}

type RegisterUserCommandType int8

const (
	UpdateUser RegisterUserCommandType = iota
	RollbackUser
	UpdateAdditional
	CancelRegistration
	RollbackAdditional
	UpdateConnectionNode
	RollbackConnectionNode
	UpdateAuth
	RollbackAuth
	ApproveRegistration
	UnknownCommand
)

type RegisterUserCommand struct {
	User User
	Type RegisterUserCommandType
}

type RegisterUserReplyType int8

const (
	UserServiceUpdated RegisterUserReplyType = iota
	UserServiceNotUpdated
	UserServiceRolledBack
	AdditionalServiceUpdated
	AdditionalServiceNotUpdated
	AdditionalServiceRolledBack
	ConnectionsUpdated
	ConnectionsNotUpdated
	ConnectionsRolledBack
	AuthUpdated
	AuthNotUpdated
	AuthRolledBack
	RegistrationCancelled
	RegistrationApproved
	UnknownReply
)

type RegisterUserReply struct {
	User User
	Type RegisterUserReplyType
}
