package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dislinkt/api_gateway/domain"
	"github.com/dislinkt/api_gateway/infrastructure/services"
	"github.com/dislinkt/api_gateway/startup/config"
	postGw "github.com/dislinkt/common/proto/post_service"
	userGw "github.com/dislinkt/common/proto/user_service"

	// "github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type PublicUserFeedHandler struct {
	postClientAddress string
	userClientAddress string
}

func NewPublicUserFeedHandler(c *config.Config) *PublicUserFeedHandler {
	return &PublicUserFeedHandler{
		postClientAddress: fmt.Sprintf("%s:%s", c.PostHost, c.PostPort),
		userClientAddress: fmt.Sprintf("%s:%s", c.UserHost, c.UserPort),
	}
}

func (handler *PublicUserFeedHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/user/public-feed", handler.GetPublicUserFeed)
	if err != nil {
		panic(err)
	}
}

func (handler *PublicUserFeedHandler) GetPublicUserFeed(w http.ResponseWriter, r *http.Request,
	pathParams map[string]string) {
	// func (handler *UserFeedHandler) GetUserFeed(ctx *gin.Context) {
	// id := ctx.Param("id")

	postClient := services.NewPostClient(handler.postClientAddress)
	userClient := services.NewUserClient(handler.userClientAddress)
	users, err := userClient.GetPublicUsers(context.TODO(), &userGw.GetMeMessage{})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	var feed []domain.Post
	for _, user := range users.Users {
		postsResponse, err := postClient.GetAllByUserId(context.TODO(), &postGw.GetRequest{Id: user.Id})
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
