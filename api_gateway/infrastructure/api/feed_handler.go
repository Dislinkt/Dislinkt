package api

import (
	"context"
	"encoding/json"
	"github.com/dislinkt/api_gateway/domain"
	"github.com/dislinkt/api_gateway/infrastructure/services"
	connectionGw "github.com/dislinkt/common/proto/connection_service"
	postGw "github.com/dislinkt/common/proto/post_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type UserFeedHandler struct {
	postClientAddress       string
	connectionClientAddress string
}

func NewUserFeedHandler(postClientAddress, connectionClientAddress string) Handler {
	return &UserFeedHandler{
		postClientAddress:       postClientAddress,
		connectionClientAddress: connectionClientAddress,
	}
}

func (handler *UserFeedHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/user/{userId}/feed", handler.GetUserFeed)
	if err != nil {
		panic(err)
	}
}

func (handler *UserFeedHandler) GetUserFeed(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := pathParams["userId"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postClient := services.NewPostClient(handler.postClientAddress)
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	connections, err := connectionClient.GetAllConnectionForUser(context.TODO(), &connectionGw.GetConnectionRequest{Uuid: id})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var feed []domain.Post
	for _, user := range connections.Users {
		postsResponse, err := postClient.GetAllByUserId(context.TODO(), &postGw.GetRequest{Id: user.UserID})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		feed = append(feed, loadUserPosts(postsResponse.Posts)...)
	}

	response, err := json.Marshal(feed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func loadUserPosts(postsPb []*postGw.Post) []domain.Post {
	var posts []domain.Post

	for _, postPb := range postsPb {
		var post domain.Post
		post.UserId = postPb.UserId
		post.PostText = postPb.PostText
		post.ImagePaths = postPb.ImagePaths
		post.Links = postPb.Links
		post.DatePosted = postPb.DatePosted
		post.Reactions = loadPostReactions(postPb.Reactions)
		post.Comments = loadPostComments(postPb.Comments)

		posts = append(posts, post)
	}
	return posts
}

func loadPostReactions(reactionsPb []*postGw.Reaction) []domain.Reaction {
	var reactions []domain.Reaction

	for _, commentPb := range reactionsPb {
		var reaction domain.Reaction
		reaction.Username = commentPb.Username
		reaction.Reaction = int(commentPb.Reaction)

		reactions = append(reactions, reaction)
	}
	return reactions
}

func loadPostComments(commentsPb []*postGw.Comment) []domain.Comment {
	var comments []domain.Comment

	for _, commentPb := range commentsPb {
		var comment domain.Comment
		comment.Username = commentPb.Username
		comment.CommentText = commentPb.CommentText

		comments = append(comments, comment)
	}
	return comments
}
