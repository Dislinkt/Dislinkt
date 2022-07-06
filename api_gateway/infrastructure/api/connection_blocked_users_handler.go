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

type ConnectionBlockedUsersHandler struct {
	userClientAddress       string
	connectionClientAddress string
}

func NewConnectionBlockedUsersHandler(c *config.Config) *ConnectionBlockedUsersHandler {
	return &ConnectionBlockedUsersHandler{
		userClientAddress:       fmt.Sprintf("%s:%s", c.UserHost, c.UserPort),
		connectionClientAddress: fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort),
	}
}

func (handler *ConnectionBlockedUsersHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/connection/user/block/{userId}", handler.GetBlockedUsers)
	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionBlockedUsersHandler) GetBlockedUsers(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	// func (handler *UserFeedHandler) GetUserFeed(ctx *gin.Context) {
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
	connections, err := connectionClient.GetAllBlockedForCurrentUser(context.TODO(),
		&connectionGw.BlockUserRequest{Uuid: id, Uuid1: ""})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	var connectionUsers []domain.ConnectionUser
	fmt.Println(connectionUsers)
	for _, user := range connections.Users {
		userResponse, err := userClient.GetOne(context.TODO(), &userGw.GetOneMessage{Id: user.UserID})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
			// ctx.AbortWithError(http.StatusBadRequest, err)
			// return
		}
		connectionUsers = append(connectionUsers, loadUserInfo(userResponse.User))
	}

	// response := feed
	response, err := json.Marshal(connectionUsers)
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

// func loadUserInfo(userPb *userGw.User) domain.ConnectionUser {
// 	var request domain.ConnectionUser
//
// 	request.UserId = userPb.Id
// 	request.Name = userPb.Name
// 	request.Surname = userPb.Surname
// 	request.Biography = userPb.Biography
// 	request.Username = userPb.Username
//
// 	return request
// }
