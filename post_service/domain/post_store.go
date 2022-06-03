package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostStore interface {
	Get(id primitive.ObjectID) (*Post, error)
	GetAll() ([]*Post, error)
	Insert(post *Post) error
	GetAllByUserId(uuid string) ([]*Post, error)
	GetAllByConnectionIds(uuids []string) ([]*Post, error)
	CreateComment(post *Post, comment *Comment) error
	LikePost(post *Post, username string) error
	DislikePost(post *Post, username string) error
	GetRecent(uuid string) ([]*Post, error)

	GetAllJobOffers() ([]*JobOffer, error)
	InsertJobOffer(offer *JobOffer) error
}
