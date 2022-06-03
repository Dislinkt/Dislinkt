package api

import (
	b64 "encoding/base64"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	pb "github.com/dislinkt/common/proto/post_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/domain"
)

func mapPost(post *domain.Post) *pb.Post {
	id := post.Id.Hex()

	links := &pb.Links{
		Comment: "/post/" + id + "/comment",
		Like:    "/post/" + id + "/like",
		Dislike: "/post/" + id + "/dislike",
		User:    "/user/" + post.UserId,
	}

	likesNum, dislikesNum := findNumberOfReactions(post)

	postPb := &pb.Post{
		Id:             id,
		UserId:         post.UserId,
		PostText:       post.PostText,
		DatePosted:     post.DatePosted.String(),
		Links:          links,
		LikesNumber:    int32(likesNum),
		DislikesNumber: int32(dislikesNum),
		CommentsNumber: int32(len(post.Comments)),
	}
	postPb.ImagePaths = convertByteToBase64(post.ImagePaths)

	return postPb
}

func findNumberOfReactions(post *domain.Post) (int, int) {
	likesNum := 0
	dislikesNum := 0

	for _, reaction := range post.Reactions {
		if reaction.Reaction == domain.LIKED {
			likesNum++
		} else if reaction.Reaction == domain.DISLIKED {
			dislikesNum++
		}
	}
	return likesNum, dislikesNum
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
		UserId:      commentPb.UserId,
		CommentText: commentPb.CommentText,
	}

	return comment
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

func mapUserReaction(user *domain.User) *pb.User {

	userPb := &pb.User{
		Username: user.Username,
		Name:     user.Name,
		Surname:  user.Surname,
	}

	return userPb
}

func mapUserCommentsForPost(user *domain.User, commentText string) *pb.Comment {

	commentPb := &pb.Comment{
		UserId:      user.UserUUID,
		Username:    user.Username,
		Name:        user.Name,
		Surname:     user.Surname,
		CommentText: commentText,
	}

	return commentPb
}

/* JOB OFFERS */

func mapJobOffer(offer *domain.JobOffer) *pb.JobOffer {
	id := offer.Id.Hex()

	offerPb := &pb.JobOffer{
		Id:            id,
		Position:      offer.Position,
		Description:   offer.Description,
		Preconditions: offer.Preconditions,
		DatePosted:    timestamppb.New(offer.DatePosted),
		Duration:      offer.Duration.String(),
		Location:      offer.Location,
		Title:         offer.Title,
		Field:         offer.Field,
	}

	return offerPb
}

func mapNewJobOffer(offerPb *pb.JobOffer) *domain.JobOffer {
	duration, _ := time.ParseDuration(offerPb.Duration)

	offer := &domain.JobOffer{
		Id:            primitive.NewObjectID(),
		Position:      offerPb.Position,
		Description:   offerPb.Description,
		Preconditions: offerPb.Preconditions,
		DatePosted:    offerPb.DatePosted.AsTime(),
		Duration:      duration,
		Location:      offerPb.Location,
		Title:         offerPb.Title,
		Field:         offerPb.Field,
	}

	return offer
}
