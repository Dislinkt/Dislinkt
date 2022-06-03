package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dislinkt/api_gateway/domain"
	"github.com/dislinkt/api_gateway/infrastructure/services"
	"github.com/dislinkt/api_gateway/startup/config"
	connectionGw "github.com/dislinkt/common/proto/connection_service"
	postGw "github.com/dislinkt/common/proto/post_service"
	// "github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type UserFeedHandler struct {
	postClientAddress       string
	connectionClientAddress string
}

func NewUserFeedHandler(c *config.Config) *UserFeedHandler {
	return &UserFeedHandler{
		postClientAddress:       fmt.Sprintf("%s:%s", c.PostHost, c.PostPort),
		connectionClientAddress: fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort),
	}
}

func (handler *UserFeedHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/user/{userId}/feed", handler.GetUserFeed)
	if err != nil {
		panic(err)
	}
}

func (handler *UserFeedHandler) GetUserFeed(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	// func (handler *UserFeedHandler) GetUserFeed(ctx *gin.Context) {
	id := pathParams["userId"]
	// id := ctx.Param("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
		// ctx.AbortWithError(http.StatusBadGateway, _)
		// return
	}

	postClient := services.NewPostClient(handler.postClientAddress)
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	connections, err := connectionClient.GetAllConnectionForUser(context.TODO(), &connectionGw.GetConnectionRequest{Uuid: id})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	var feed []domain.Post
	for _, user := range connections.Users {
		postsResponse, err := postClient.GetAllByUserId(context.TODO(), &postGw.GetRequest{Id: user.UserID})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
			// ctx.AbortWithError(http.StatusBadRequest, err)
			// return
		}
		feed = append(feed, loadUserPosts(postsResponse.Posts)...)
	}

	// response := feed
	response, err := json.Marshal(feed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
		// ctx.AbortWithError(http.StatusInternalServerError, err)
		// return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	// ctx.JSON(http.StatusOK, &response)
}

func loadUserPosts(postsPb []*postGw.Post) []domain.Post {
	var posts []domain.Post

	for _, postPb := range postsPb {
		var post domain.Post
		post.Id = postPb.Id
		post.UserId = postPb.UserId
		post.PostText = postPb.PostText
		post.ImagePaths = postPb.ImagePaths
		post.DatePosted = postPb.DatePosted
		post.LikesNumber = int(postPb.LikesNumber)
		post.DislikesNumber = int(postPb.DislikesNumber)
		post.CommentsNumber = int(postPb.CommentsNumber)
		post.Links = domain.Links{Comment: postPb.Links.Comment, Dislike: postPb.Links.Dislike, Like: postPb.Links.Like}

		posts = append(posts, post)
	}
	return posts
}
