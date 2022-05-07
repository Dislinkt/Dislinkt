package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	Id         primitive.ObjectID `bson:"_id"`
	UserId     string             `bson:"user_id"`
	PostText   string             `bson:"post_text"`
	ImagePaths []string           `bson:"image_paths"`
	Links      []string           `bson:"links"`
	DatePosted time.Time          `bson:"date_posted"`
	Reactions  []Reaction         `bson:"reactions"`
	Comments   []Comment          `bson:"comments"`
	IsDeleted  bool               `bson:"is_deleted"`
}

type Comment struct {
	Username    string `bson:"username"`
	CommentText string `bson:"comment_text"`
}

type Reaction struct {
	Username string       `bson:"username"`
	Reaction ReactionType `bson:"reaction"`
}

type ReactionType int

const (
	Neutral ReactionType = iota
	LIKED
	DISLIKED
)
