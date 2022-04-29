package api

import (
	pb "github.com/dislinkt/common/proto/post_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/domain"
	"time"
)

func mapPost(post *domain.Post) *pb.Post {
	postPb := &pb.Post{
		Id:         post.Id.Hex(),
		UserId:     post.UserId,
		PostText:   post.PostText,
		DatePosted: post.DatePosted.String(),
	}
	for _, image := range post.ImagePaths {
		postPb.ImagePaths = append(postPb.ImagePaths, image)
	}
	for _, link := range post.Links {
		postPb.Links = append(postPb.Links, link)
	}
	for _, reaction := range post.Reactions {
		postPb.Reactions = append(postPb.Reactions, &pb.Reaction{
			Username: reaction.Username,
			Reaction: int32(reaction.Reaction),
		})
	}
	for _, comment := range post.Comments {
		postPb.Comments = append(postPb.Comments, &pb.Comment{
			Username:    comment.Username,
			CommentText: comment.CommentText,
		})
	}
	return postPb
}

func mapNewPost(postPb *pb.Post) *domain.Post {
	post := &domain.Post{
		Id:         primitive.NewObjectID(),
		UserId:     postPb.UserId,
		PostText:   postPb.PostText,
		DatePosted: time.Now(),
	}
	for _, image := range postPb.ImagePaths {
		post.ImagePaths = append(post.ImagePaths, image)
	}
	for _, link := range postPb.Links {
		post.Links = append(post.Links, link)
	}

	return post
}
