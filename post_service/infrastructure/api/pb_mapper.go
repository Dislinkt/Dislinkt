package api

import (
	b64 "encoding/base64"
	"fmt"
	"time"

	pb "github.com/dislinkt/common/proto/post_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/domain"
)

func mapPost(post *domain.Post) *pb.Post {
	postPb := &pb.Post{
		Id:         post.Id.Hex(),
		UserId:     post.UserId,
		PostText:   post.PostText,
		DatePosted: post.DatePosted.String(),
	}
	postPb.ImagePaths = convertByteToBase64(post.ImagePaths)
	for _, reaction := range post.Reactions {
		postPb.Reactions = append(postPb.Reactions, &pb.Reaction{
			Username: reaction.Username,
			Reaction: mapReactionTypeToPb(reaction.Reaction),
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
	post.ImagePaths = convertBase64ToByte(postPb.ImagePaths)

	return post
}

func mapNewComment(commentPb *pb.Comment) *domain.Comment {
	comment := &domain.Comment{
		Username:    commentPb.Username,
		CommentText: commentPb.CommentText,
	}

	return comment
}

func mapReactionTypeToPb(reactionType domain.ReactionType) pb.ReactionType {
	switch reactionType {
	case domain.Neutral:
		return pb.ReactionType_Neutral
	case domain.LIKED:
		return pb.ReactionType_LIKED
	case domain.DISLIKED:
		return pb.ReactionType_DISLIKED
	}
	return pb.ReactionType_Neutral
}

func convertBase64ToByte(images []string) [][]byte {
	var decodedImages [][]byte
	for _, image := range images {
		fmt.Println(image)
		imageDec, _ := b64.StdEncoding.DecodeString(image)
		fmt.Println(string(imageDec))
		decodedImages = append(decodedImages, imageDec)
	}
	return decodedImages
}
func convertByteToBase64(images [][]byte) []string {
	var encodedImages []string
	for _, image := range images {
		fmt.Println(image)
		imageEnc := b64.StdEncoding.EncodeToString(image)
		fmt.Println(string(imageEnc))
		encodedImages = append(encodedImages, imageEnc)
	}
	return encodedImages
}
