package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"post_service/domain"
)

const (
	DATABASE             = "post"
	COLLECTION_POST      = "post"
	COLLECTION_JOB_OFFER = "job_offer"
	COLLECTION_USER      = "user"
)

type PostMongoDBStore struct {
	posts     *mongo.Collection
	jobOffers *mongo.Collection
	users     *mongo.Collection
}

func NewPostMongoDBStore(client *mongo.Client) domain.PostStore {
	posts := client.Database(DATABASE).Collection(COLLECTION_POST)
	jobOffers := client.Database(DATABASE).Collection(COLLECTION_JOB_OFFER)
	users := client.Database(DATABASE).Collection(COLLECTION_USER)
	return &PostMongoDBStore{
		posts:     posts,
		jobOffers: jobOffers,
		users:     users,
	}
}

func (store *PostMongoDBStore) GetRecent(uuid string) ([]*domain.Post, error) {

	filter := bson.M{"user_id": uuid, "date_posted": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Date(time.Now().Year(), time.Now().Month()-1, time.Now().Day(), 0, 0, 0, 0, &time.Location{}))}}
	return store.filter(filter)
}

func (store *PostMongoDBStore) Get(id primitive.ObjectID) (*domain.Post, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *PostMongoDBStore) GetAll() ([]*domain.Post, error) {
	filter := bson.D{}
	return store.filter(filter)
}

func (store *PostMongoDBStore) GetAllByUserId(uuid string) ([]*domain.Post, error) {
	filter := bson.M{"user_id": uuid}
	return store.filter(filter)
}

func (store *PostMongoDBStore) GetAllByConnectionIds(uuids []string) ([]*domain.Post, error) {
	var posts []*domain.Post

	for _, uuid := range uuids {
		postsByUser, err := store.GetAllByUserId(uuid)
		posts = append(posts, postsByUser...)

		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (store *PostMongoDBStore) Insert(post *domain.Post) error {
	result, err := store.posts.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	post.Id = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (store *PostMongoDBStore) CreateComment(post *domain.Post, comment *domain.Comment) error {
	comments := append(post.Comments, *comment)

	_, err := store.posts.UpdateOne(context.TODO(), bson.M{"_id": post.Id}, bson.D{
		{"$set", bson.D{{"comments", comments}}},
	},
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostMongoDBStore) LikePost(post *domain.Post, userId string) error {

	var reactions []domain.Reaction

	reactionExists := false
	for _, reaction := range post.Reactions {
		if reaction.UserId != userId {
			reactions = append(reactions, reaction)
		} else {
			if reaction.Reaction != domain.LIKED {
				reaction.Reaction = domain.LIKED
				reactions = append(reactions, reaction)
			}
			reactionExists = true
		}

	}
	if !reactionExists {
		reaction := domain.Reaction{
			UserId:   userId,
			Reaction: domain.LIKED,
		}
		reactions = append(reactions, reaction)
	}

	_, err := store.posts.UpdateOne(context.TODO(), bson.M{"_id": post.Id}, bson.D{
		{"$set", bson.D{{"reactions", reactions}}},
	},
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostMongoDBStore) DislikePost(post *domain.Post, userId string) error {
	var reactions []domain.Reaction

	reactionExists := false
	for _, reaction := range post.Reactions {
		if reaction.UserId != userId {
			reactions = append(reactions, reaction)
		} else {
			if reaction.Reaction != domain.DISLIKED {
				reaction.Reaction = domain.DISLIKED
				reactions = append(reactions, reaction)
			}
			reactionExists = true
		}

	}
	if !reactionExists {
		reaction := domain.Reaction{
			UserId:   userId,
			Reaction: domain.DISLIKED,
		}
		reactions = append(reactions, reaction)
	}

	_, err := store.posts.UpdateOne(context.TODO(), bson.M{"_id": post.Id}, bson.D{
		{"$set", bson.D{{"reactions", reactions}}},
	},
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostMongoDBStore) filter(filter interface{}) ([]*domain.Post, error) {
	cursor, err := store.posts.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	return decode(cursor)
}

func (store *PostMongoDBStore) filterOne(filter interface{}) (post *domain.Post, err error) {
	result := store.posts.FindOne(context.TODO(), filter)
	err = result.Decode(&post)
	return
}

func decode(cursor *mongo.Cursor) (posts []*domain.Post, err error) {
	for cursor.Next(context.TODO()) {
		var post domain.Post
		err = cursor.Decode(&post)
		if err != nil {
			return
		}
		posts = append(posts, &post)
	}
	err = cursor.Err()
	return
}

/* JOB OFFERS */

func (store *PostMongoDBStore) InsertJobOffer(offer *domain.JobOffer) error {
	result, err := store.jobOffers.InsertOne(context.TODO(), offer)
	fmt.Println("PostMongoDBStore: InsertJobOffer")
	if err != nil {
		return err
	}
	fmt.Println(result)
	offer.Id = result.InsertedID.(primitive.ObjectID)
	fmt.Println(offer.Id)
	return nil
}

func (store *PostMongoDBStore) DeleteJobOffer(jobOffer *domain.JobOffer) error {
	_, err := store.users.DeleteOne(context.TODO(), bson.M{"_id": jobOffer.Id})
	return err
}

func (store *PostMongoDBStore) GetAllJobOffers() ([]*domain.JobOffer, error) {
	filter := bson.D{}
	return store.filterJobOffers(filter)
}

func (store *PostMongoDBStore) SearchJobOffers(searchText string) ([]*domain.JobOffer, error) {
	filter := bson.M{
		"$or": []bson.M{
			{
				"title": bson.M{
					"$regex": primitive.Regex{
						Pattern: searchText,
						Options: "i",
					},
				},
			},
			{
				"position": bson.M{
					"$regex": primitive.Regex{
						Pattern: searchText,
						Options: "i",
					},
				},
			},
			{
				"location": bson.M{
					"$regex": primitive.Regex{
						Pattern: searchText,
						Options: "i",
					},
				},
			},
			{
				"field": bson.M{
					"$regex": primitive.Regex{
						Pattern: searchText,
						Options: "i",
					},
				},
			},
		},
	}
	return store.filterJobOffers(filter)
}

func (store *PostMongoDBStore) filterJobOffers(filter interface{}) ([]*domain.JobOffer, error) {
	cursor, err := store.jobOffers.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	return decodeJobOffers(cursor)
}

func decodeJobOffers(cursor *mongo.Cursor) (offers []*domain.JobOffer, err error) {
	for cursor.Next(context.TODO()) {
		var offer domain.JobOffer
		err = cursor.Decode(&offer)
		if err != nil {
			return
		}
		offers = append(offers, &offer)
	}
	err = cursor.Err()
	return
}

/* USERS */

func (store *PostMongoDBStore) InsertUser(user *domain.User) error {
	user.Id = primitive.NewObjectID()
	_, err := store.users.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostMongoDBStore) DeleteUser(user *domain.User) error {
	_, err := store.users.DeleteOne(context.TODO(), bson.M{"userUUID": user.UserUUID})
	return err
}

func (store *PostMongoDBStore) UpdateUser(user *domain.User) error {

	_, err := store.users.UpdateOne(context.TODO(), bson.M{"userUUID": user.UserUUID}, bson.D{
		{"$set", bson.D{{"username", user.Username}}},
		{"$set", bson.D{{"name", user.Name}}},
		{"$set", bson.D{{"surname", user.Surname}}},
	},
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostMongoDBStore) GetUser(id string) (*domain.User, error) {
	filter := bson.M{"userUUID": id}
	return store.filterOneUser(filter)
}

func (store *PostMongoDBStore) filterOneUser(filter interface{}) (user *domain.User, err error) {
	result := store.users.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}
