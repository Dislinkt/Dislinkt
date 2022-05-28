package domain

type AdditionalUserStore interface {
	FindUserDocument(uuid string) (*AdditionalUser, error)
	CreateUserDocument(uuid string) (*AdditionalUser, error)
	DeleteUserDocument(uuid string) error
	FindDocument(uuid string) (*AdditionalUser, error)

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

	// ENTITES
	InsertFieldOfStudy(filedOfStudies []*FieldOfStudy) ([]*FieldOfStudy, error)
	GetAllFieldOfStudy() ([]*FieldOfStudy, error)
	InsertSkills(skills []*Skill) ([]*Skill, error)
	GetSkills() (skills []*Skill, err error)
	InsertIndustries(industries []*Industry) ([]*Industry, error)
	GetIndustries() (industries []*Industry, err error)
	InsertDegrees(degrees []*Degree) ([]*Degree, error)
	GetDegrees() (degrees []*Degree, err error)
}
