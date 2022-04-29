package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdditionalUser struct {
	Id         primitive.ObjectID `bson:"_id"`
	UserUUID   string             `bson:"userUUID"`
	Educations []Education        `bson:"educations"`
	Positions  []Position         `bson:"positions"`
	Skills     []Skill            `bson:"skills"`
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
