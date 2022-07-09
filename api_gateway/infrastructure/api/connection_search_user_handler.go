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
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type ConnectionSearchUserHandler struct {
	userClientAddress       string
	connectionClientAddress string
}

func NewConnectionSearchUserHandler(c *config.Config) *ConnectionSearchUserHandler {
	return &ConnectionSearchUserHandler{
		userClientAddress:       fmt.Sprintf("%s:%s", c.UserHost, c.UserPort),
		connectionClientAddress: fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort),
	}
}

func (handler *ConnectionSearchUserHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/search/{userId}", handler.GetConnectionRequests)
	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionSearchUserHandler) GetConnectionRequests(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	// func (handler *UserFeedHandler) GetUserFeed(ctx *gin.Context) {
	searchText := r.URL.Query().Get("searchText")
	id := pathParams["userId"]
	// id := ctx.Param("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
		// ctx.AbortWithError(http.StatusBadGateway, _)
		// return
	}
	userClient := services.NewUserClient(handler.userClientAddress)
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	users, err := userClient.GetAll(context.TODO(), &userGw.SearchMessage{SearchText: searchText})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	var allowedUsers []domain.ConnectionUser
	for _, user := range users.Users {
		userResponse, err := connectionClient.CheckIfUsersBlocked(context.TODO(), &connectionGw.CheckConnection{
			Uuid1: id,
			Uuid2: user.Id,
		})
		userResponse2, err := connectionClient.CheckIfUsersBlocked(context.TODO(), &connectionGw.CheckConnection{
			Uuid1: user.Id,
			Uuid2: id,
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
			// ctx.AbortWithError(http.StatusBadRequest, err)
			// return
		}
		if !userResponse.IsBlocked && !userResponse2.IsBlocked {
			allowedUsers = append(allowedUsers, loadUserInfo(user))
		}
	}

	// response := feed
	response, err := json.Marshal(allowedUsers)
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
