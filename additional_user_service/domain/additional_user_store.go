package domain

type AdditionalUserStore interface {
	FindUserDocument(uuid string) (*AdditionalUser, error)
	CreateUserDocument(uuid string) (*AdditionalUser, error)
	FindOrCreateDocument(uuid string) (*AdditionalUser, error)

	// EDUCATION
	InsertEducation(uuid string, education *Education) (*Education, error)
	// GetUserEducation(educationId string) (*Education, error)
	UpdateUserEducation(educationId string, education *Education) error
	DeleteUserEducation(id string) error

	// POSITION
	UpdateUserPosition(uuid string, position *Position) error
	InsertPosition(educationId string, position *Position) (*Position, error)
	DeleteUserPosition(id string) error

	// SKILL
	InsertSkill(uuid string, skill *Skill) (*Skill, error)
	UpdateUserSkill(id string, skill *Skill) error
	DeleteUserSkill(id string) error

	// INTEREST
	InsertInterest(uuid string, interest *Interest) (*Interest, error)
	UpdateUserInterest(id string, interest *Interest) error
	DeleteUserInterest(id string) error
}
