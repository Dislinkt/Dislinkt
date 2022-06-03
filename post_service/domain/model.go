package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Id         primitive.ObjectID `bson:"_id"`
	UserId     string             `bson:"user_id"`
	PostText   string             `bson:"post_text"`
	ImagePaths [][]byte           `bson:"image_paths"`
	DatePosted time.Time          `bson:"date_posted"`
	Reactions  []Reaction         `bson:"reactions"`
	Comments   []Comment          `bson:"comments"`
	IsDeleted  bool               `bson:"is_deleted"`
}

type Comment struct {
	UserId      string `bson:"user_id"`
	CommentText string `bson:"comment_text"`
}

type Reaction struct {
	UserId   string       `bson:"user_id"`
	Reaction ReactionType `bson:"reaction"`
}

type ReactionType int

const (
	Neutral ReactionType = iota
	LIKED
	DISLIKED
)

type JobOffer struct {
	Id            primitive.ObjectID `bson:"_id"`
	Position      string             `bson:"position"`
	Description   string             `bson:"description"`
	Preconditions string             `bson:"preconditions"`
	DatePosted    time.Time          `bson:"date_posted"`
	Duration      time.Duration      `bson:"duration"`
	Location      string             `bson:"location"`
	Title         string             `bson:"title"`
	Field         string             `bson:"field"`
}

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	UserUUID string             `bson:"userUUID"`
	Username string             `bson:"username"`
	Name     string             `bson:"name"`
	Surname  string             `bson:"surname"`
}
