package persistence

import (
	"context"
	"fmt"

	"github.com/dislinkt/additional_user_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "additional_user"
	COLLECTION = "user"
)

type AdditionalUserMongoDBStore struct {
	users *mongo.Collection
}

func NewAdditionalUserMongoDBStore(client *mongo.Client) domain.AdditionalUserStore {
	users := client.Database(DATABASE).Collection(COLLECTION)
	return &AdditionalUserMongoDBStore{
		users: users,
	}
}

func (store *AdditionalUserMongoDBStore) InsertEducation(uuid string, education *domain.Education) (*domain.Education, error) {
	education.Id = primitive.NewObjectID()
	_, err := store.users.UpdateOne(context.TODO(), bson.M{"userUUID": uuid}, bson.D{
		{"$push", bson.D{{"educations", education}}},
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return education, nil
}

func (store *AdditionalUserMongoDBStore) FindUserDocument(userUUID string) (user *domain.AdditionalUser, err error) {
	err = store.users.FindOne(context.TODO(), bson.D{{"userUUID", userUUID}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (store *AdditionalUserMongoDBStore) CreateUserDocument(uuid string) (*domain.AdditionalUser, error) {
	userDocumentEmpty := domain.AdditionalUserEmpty{
		Id:       primitive.NewObjectID(),
		UserUUID: uuid,
	}
	userDocument := domain.AdditionalUser{
		Id:       primitive.NewObjectID(),
		UserUUID: uuid,
	}
	result, err := store.users.InsertOne(context.TODO(), userDocumentEmpty)
	if err != nil {
		return nil, err
	}
	userDocument.Id = result.InsertedID.(primitive.ObjectID)
	return &userDocument, nil
}

func (store *AdditionalUserMongoDBStore) FindOrCreateDocument(uuid string) (*domain.AdditionalUser, error) {
	userDocument, err := store.FindUserDocument(uuid)
	if err != nil {
		userDocument, err = store.CreateUserDocument(uuid)
		if err != nil {
			return nil, err
		}
	}
	return userDocument, nil
}
