package domain

type AdditionalUserStore interface {
	InsertEducation(uuid string, education *Education) (*Education, error)
	FindUserDocument(userUUID string) (*AdditionalUser, error)
	CreateUserDocument(uuid string) (*AdditionalUser, error)
	FindOrCreateDocument(uuid string) (*AdditionalUser, error)
}
