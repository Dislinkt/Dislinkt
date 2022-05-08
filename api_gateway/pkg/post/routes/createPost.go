package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/post_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Reaction struct {
	Username string `json:"username"`
	Reaction int    `json:"reaction"`
}

type Comment struct {
	Username    string `json:"username"`
	CommentText string `json:"commentText"`
}

type PostRequestBody struct {
	Id         string     `json:"id"`
	UserId     string     `json:"user_id"`
	PostText   string     `json:"post_text"`
	ImagePaths []string   `json:"image_paths"`
	Links      []string   `json:"links"`
	DatePosted string     `json:"date_posted"`
	Reactions  []Reaction `json:"reactions"`
	Comments   []Comment  `json:"comments"`
}

func CreatePost(ctx *gin.Context, c pb.PostServiceClient) {
	b := PostRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.CreatePost(context.Background(), &pb.CreatePostRequest{
		Post: &pb.Post{
			Id:         b.Id,
			UserId:     b.UserId,
			PostText:   b.PostText,
			ImagePaths: b.ImagePaths,
			Links:      b.Links,
			DatePosted: b.DatePosted,
			Reactions:  loadPostReactions(b.Reactions),
			Comments:   loadComments(b.Comments),
		},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}

func loadPostReactions(reactions []Reaction) []*pb.Reaction {
	var reactionsPb []*pb.Reaction

	for _, reaction := range reactions {
		var reactionPb *pb.Reaction
		reactionPb.Username = reaction.Username
		reactionPb.Reaction = pb.ReactionType(reaction.Reaction)

		reactionsPb = append(reactionsPb, reactionPb)
	}
	return reactionsPb
}

func loadComments(comments []Comment) []*pb.Comment {
	var commentsPb []*pb.Comment

	for _, comment := range comments {
		var commentPb *pb.Comment
		commentPb.Username = comment.Username
		commentPb.CommentText = comment.CommentText

		commentsPb = append(commentsPb, commentPb)
	}
	return commentsPb
}
