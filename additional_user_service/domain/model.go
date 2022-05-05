package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdditionalUser struct {
	Id         primitive.ObjectID    `bson:"_id"`
	UserUUID   string                `bson:"userUUID"`
	Educations *map[string]Education `bson:"educations,omitempty"`
	Positions  *map[string]Position  `bson:"positions,omitempty"`
	Skills     *map[string]Skill     `bson:"skills,omitempty"`
	Interests  *map[string]Interest  `bson:"interests,omitempty"`
}

type AdditionalUserEmpty struct {
	Id       primitive.ObjectID `bson:"_id"`
	UserUUID string             `bson:"userUUID"`
}
type Education struct {
	Id           primitive.ObjectID `bson:"_id"`
	School       string             `bson:"school"`
	Degree       string             `bson:"degree"`
	FieldOfStudy string             `bson:"field_of_study"`
	StartDate    primitive.DateTime `bson:"start_date"`
	EndDate      primitive.DateTime `bson:"end_date"`
}

type Position struct {
	Id          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	CompanyName string             `bson:"company_name"`
	Industry    string             `bson:"industry"`
	StartDate   primitive.DateTime `bson:"start_date"`
	EndDate     primitive.DateTime `bson:"end_date"`
	Current     bool               `bson:"current"`
}

type Skill struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type Interest struct {
	Id    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Group InterestGroup      `bson:"group"`
}

type InterestGroup string

const (
	Group1 InterestGroup = "GROUP_1"
	Group2               = "GROUP_2"
	Group3               = "GROUP_3"
)
