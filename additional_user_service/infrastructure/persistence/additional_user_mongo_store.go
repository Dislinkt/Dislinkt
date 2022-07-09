package persistence

import (
	"context"
	"fmt"
	"log"

	"github.com/dislinkt/additional_user_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DATABASE                  = "additional_user"
	COLLECTION_USER           = "user"
	COLLECTION_FIELD_OF_STUDY = "field_of_study"
	COLLECTION_SKILLS         = "skills"
	COLLECTION_INDUSTRIES     = "industries"
	COLLECTION_DEGREES        = "degrees"
)

type AdditionalUserMongoDBStore struct {
	users        *mongo.Collection
	fieldOfStudy *mongo.Collection
	skills       *mongo.Collection
	industries   *mongo.Collection
	degrees      *mongo.Collection
}

func NewAdditionalUserMongoDBStore(client *mongo.Client) domain.AdditionalUserStore {
	users := client.Database(DATABASE).Collection(COLLECTION_USER)
	fieldOfStudy := client.Database(DATABASE).Collection(COLLECTION_FIELD_OF_STUDY)
	skills := client.Database(DATABASE).Collection(COLLECTION_SKILLS)
	industries := client.Database(DATABASE).Collection(COLLECTION_INDUSTRIES)
	degrees := client.Database(DATABASE).Collection(COLLECTION_DEGREES)
	return &AdditionalUserMongoDBStore{
		users:        users,
		fieldOfStudy: fieldOfStudy,
		skills:       skills,
		industries:   industries,
		degrees:      degrees,
	}
}

func (store *AdditionalUserMongoDBStore) InsertDegrees(degrees []*domain.Degree) ([]*domain.Degree, error) {
	fmt.Println("InsertDegree")
	fmt.Println(degrees)
	for _, i := range degrees {
		_, err := store.degrees.InsertOne(context.TODO(), i)
		if err != nil {
			return nil, err
		}
	}

	return degrees, nil
}

func (store *AdditionalUserMongoDBStore) GetDegrees() (degrees []*domain.Degree, err error) {
	cursor, err := store.degrees.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var allDegrees []*domain.Degree
	for cursor.Next(context.TODO()) {
		var degree domain.Degree
		if err = cursor.Decode(&degree); err != nil {
			log.Fatal(err)
		}
		fmt.Println(degree)
		allDegrees = append(allDegrees, &degree)
	}
	fmt.Println(allDegrees)
	return allDegrees, nil
}

func (store *AdditionalUserMongoDBStore) InsertIndustries(industries []*domain.Industry) ([]*domain.Industry, error) {
	fmt.Println("InsertIndustries")
	fmt.Println(industries)
	for _, i := range industries {
		_, err := store.industries.InsertOne(context.TODO(), i)
		if err != nil {
			return nil, err
		}
	}

	return industries, nil
}

func (store *AdditionalUserMongoDBStore) GetIndustries() (industries []*domain.Industry, err error) {
	cursor, err := store.industries.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var allIndustries []*domain.Industry
	for cursor.Next(context.TODO()) {
		var industry domain.Industry
		if err = cursor.Decode(&industry); err != nil {
			log.Fatal(err)
		}
		fmt.Println(industry)
		allIndustries = append(allIndustries, &industry)
	}
	fmt.Println(allIndustries)
	return allIndustries, nil
}

func (store *AdditionalUserMongoDBStore) InsertSkills(skills []*domain.Skill) ([]*domain.Skill, error) {
	fmt.Println("InsertSkills")
	fmt.Println(skills)
	for _, s := range skills {
		_, err := store.skills.InsertOne(context.TODO(), s)
		if err != nil {
			return nil, err
		}
	}

	return skills, nil
}

func (store *AdditionalUserMongoDBStore) InsertFieldOfStudy(filedOfStudies []*domain.FieldOfStudy) ([]*domain.
	FieldOfStudy, error) {
	fmt.Println("InsertFieldOfStudy")
	fmt.Println(filedOfStudies)
	for _, f := range filedOfStudies {
		_, err := store.fieldOfStudy.InsertOne(context.TODO(), f)
		if err != nil {
			return nil, err
		}
	}

	return filedOfStudies, nil
}
func (store *AdditionalUserMongoDBStore) GetSkills() ([]*domain.Skill, error) {
	cursor, err := store.skills.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var skills []*domain.Skill
	for cursor.Next(context.TODO()) {
		var skill domain.Skill
		if err = cursor.Decode(&skill); err != nil {
			log.Fatal(err)
		}
		fmt.Println(skill)
		skills = append(skills, &skill)
	}
	fmt.Println(skills)
	return skills, nil
}

func (store *AdditionalUserMongoDBStore) GetAllFieldOfStudy() ([]*domain.FieldOfStudy, error) {
	cursor, err := store.fieldOfStudy.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var fields []*domain.FieldOfStudy
	for cursor.Next(context.TODO()) {
		var field domain.FieldOfStudy
		if err = cursor.Decode(&field); err != nil {
			log.Fatal(err)
		}
		fmt.Println(field)
		fields = append(fields, &field)
	}
	fmt.Println(fields)
	return fields, nil
}

func (store *AdditionalUserMongoDBStore) FindUserDocument(userUUID string) (user *domain.AdditionalUser, err error) {
	err = store.users.FindOne(context.TODO(), bson.D{{"userUUID", userUUID}}).Decode(&user)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User with provided id does not exist.")
	}
	return user, nil
}

func (store *AdditionalUserMongoDBStore) CreateUserDocument(uuid string) (*domain.AdditionalUser, error) {
	userDocument := domain.AdditionalUser{
		Id:       primitive.NewObjectID(),
		UserUUID: uuid,
	}
	result, err := store.users.InsertOne(context.TODO(), userDocument)
	if err != nil {
		return nil, err
	}
	userDocument.Id = result.InsertedID.(primitive.ObjectID)
	return &userDocument, nil
}

func (store *AdditionalUserMongoDBStore) DeleteUserDocument(uuid string) error {
	_, err := store.users.DeleteOne(context.TODO(), bson.M{"userUUID": uuid})
	return err
}

func (store *AdditionalUserMongoDBStore) FindDocument(uuid string) (*domain.AdditionalUser, error) {
	userDocument, err := store.FindUserDocument(uuid)
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	return userDocument, nil
}

// EDUCATION

func (store *AdditionalUserMongoDBStore) InsertEducation(uuid string,
	education *domain.Education) (*domain.Education, error) {
	_, err := store.users.UpdateOne(context.TODO(), bson.M{"userUUID": uuid}, bson.D{
		{"$set", bson.D{{"educations." + education.Id.Hex(), education}}},
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("[AdditionalUserMongoDBStore InsertEducation]")
	fmt.Println(education)
	return education, nil
}

func (store *AdditionalUserMongoDBStore) UpdateUserEducation(educationId string,
	education *domain.Education) error {
	id, err := primitive.ObjectIDFromHex(educationId)
	if err != nil {
		return status.Error(codes.NotFound, "Education with provided id does not exist.")
	}

	update := bson.D{{"$set",
		bson.D{
			{"educations." + educationId + ".degree", education.Degree},
			{"educations." + educationId + ".school", education.School},
			{"educations." + educationId + ".field_of_study", education.FieldOfStudy},
			{"educations." + educationId + ".start_date", education.StartDate},
			{"educations." + educationId + ".end_date", education.EndDate},
		},
	}}
	_, err = store.users.UpdateOne(context.TODO(), bson.M{"educations." + educationId + "._id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// func (store *AdditionalUserMongoDBStore) GetUserEducation(educationId string) (*domain.Education, error) {
// 	id, err := primitive.ObjectIDFromHex(educationId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pipeline := mongo.Pipeline{
// 		{
// 			{"$unwind", "$educations"},
// 		},
// 		{
// 			{"$match", bson.D{
// 				{"educations." + educationId + "._id", id},
// 			}},
// 		},
// 	}
// 	cursor, err := store.users.Aggregate(context.TODO(), pipeline)
// 	var results []domain.AdditionalUser
// 	if err := cursor.All(context.TODO(), &results); err != nil {
// 		return nil, err
// 	}
// 	nesto := (*results[0].Educations)[educationId]
//
// 	return &nesto, nil
// }

func (store *AdditionalUserMongoDBStore) DeleteUserEducation(educationId string) error {
	id, err := primitive.ObjectIDFromHex(educationId)
	if err != nil {
		return status.Error(codes.NotFound, "Education with provided id does not exist.")
	}

	update := bson.D{{"$unset",
		bson.D{
			{"educations." + educationId, ""},
		},
	}}
	_, err = store.users.UpdateOne(context.TODO(), bson.M{"educations." + educationId + "._id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// POSITION

func (store *AdditionalUserMongoDBStore) UpdateUserPosition(positionId string, position *domain.Position) error {

	id, err := primitive.ObjectIDFromHex(positionId)
	if err != nil {
		return status.Error(codes.NotFound, "Position with provided id does not exist.")
	}

	update := bson.D{{"$set",
		bson.D{
			{"positions." + positionId + ".title", position.Title},
			{"positions." + positionId + ".company_name", position.CompanyName},
			{"positions." + positionId + ".industry", position.Industry},
			{"positions." + positionId + ".start_date", position.StartDate},
			{"positions." + positionId + ".end_date", position.EndDate},
			{"positions." + positionId + ".current", position.Current},
		},
	}}
	_, err = store.users.UpdateOne(context.TODO(), bson.M{"positions." + positionId + "._id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

func (store *AdditionalUserMongoDBStore) InsertPosition(uuid string, position *domain.Position) (*domain.Position, error) {
	position.Id = primitive.NewObjectID()
	_, err := store.users.UpdateOne(context.TODO(), bson.M{"userUUID": uuid}, bson.D{
		{"$set", bson.D{{"positions." + position.Id.Hex(), position}}},
	})
	if err != nil {
		return nil, err
	}
	return position, nil
}

func (store *AdditionalUserMongoDBStore) DeleteUserPosition(positionId string) error {
	id, err := primitive.ObjectIDFromHex(positionId)
	if err != nil {
		return status.Error(codes.NotFound, "Position with provided id does not exist.")
	}

	update := bson.D{{"$unset",
		bson.D{
			{"positions." + positionId, ""},
		},
	}}
	_, err = store.users.UpdateOne(context.TODO(), bson.M{"positions." + positionId + "._id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// SKILL

func (store *AdditionalUserMongoDBStore) InsertSkill(uuid string, skill *domain.Skill) (*domain.Skill, error) {
	_, err := store.users.UpdateOne(context.TODO(), bson.M{"userUUID": uuid}, bson.D{
		{"$set", bson.D{{"skills." + skill.Id.Hex(), skill}}},
	})
	if err != nil {
		return nil, err
	}
	return skill, nil
}

func (store *AdditionalUserMongoDBStore) UpdateUserSkill(skillId string, skill *domain.Skill) error {

	id, err := primitive.ObjectIDFromHex(skillId)
	if err != nil {
		return status.Error(codes.NotFound, "Skill with provided id does not exist.")
	}

	update := bson.D{{"$set",
		bson.D{
			{"skills." + skillId + ".name", skill.Name},
		},
	}}
	_, err = store.users.UpdateOne(context.TODO(), bson.M{"skills." + skillId + "._id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

func (store *AdditionalUserMongoDBStore) DeleteUserSkill(skillId string) error {
	id, err := primitive.ObjectIDFromHex(skillId)
	if err != nil {
		return status.Error(codes.NotFound, "Skill with provided id does not exist.")
	}

	update := bson.D{{"$unset",
		bson.D{
			{"skills." + skillId, ""},
		},
	}}
	_, err = store.users.UpdateOne(context.TODO(), bson.M{"skills." + skillId + "._id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// INTEREST

func (store *AdditionalUserMongoDBStore) InsertInterest(uuid string, interest *domain.Interest) (*domain.Interest, error) {
	interest.Id = primitive.NewObjectID()
	_, err := store.users.UpdateOne(context.TODO(), bson.M{"userUUID": uuid}, bson.D{
		{"$set", bson.D{{"interests." + interest.Id.Hex(), interest}}},
	})
	if err != nil {
		return nil, err
	}
	return interest, nil
}

func (store *AdditionalUserMongoDBStore) UpdateUserInterest(interestId string, interest *domain.Interest) error {

	id, err := primitive.ObjectIDFromHex(interestId)
	if err != nil {
		return status.Error(codes.NotFound, "Position with provided id does not exist.")
	}

	update := bson.D{{"$set",
		bson.D{
			{"interests." + interestId + ".name", interest.Name},
			{"interests." + interestId + ".group", interest.Group},
		},
	}}
	_, err = store.users.UpdateOne(context.TODO(), bson.M{"interests." + interestId + "._id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

func (store *AdditionalUserMongoDBStore) DeleteUserInterest(interestId string) error {
	id, err := primitive.ObjectIDFromHex(interestId)
	if err != nil {
		return status.Error(codes.NotFound, "Interest with provided id does not exist.")
	}

	update := bson.D{{"$unset",
		bson.D{
			{"interests." + interestId, ""},
		},
	}}
	_, err = store.users.UpdateOne(context.TODO(), bson.M{"interests." + interestId + "._id": id}, update)
	if err != nil {
		return err
	}

	return nil
}
