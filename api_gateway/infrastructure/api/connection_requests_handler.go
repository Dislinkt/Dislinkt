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

type ConnectionRequestHandler struct {
	userClientAddress       string
	connectionClientAddress string
}

func NewConnectionRequestHandler(c *config.Config) *ConnectionRequestHandler {
	return &ConnectionRequestHandler{
		userClientAddress:       fmt.Sprintf("%s:%s", c.UserHost, c.UserPort),
		connectionClientAddress: fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort),
	}
}

func (handler *ConnectionRequestHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/connection/user/{userId}/request", handler.GetConnectionRequests)
	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionRequestHandler) GetConnectionRequests(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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
	connections, err := connectionClient.GetAllConnectionRequestsForUser(context.TODO(), &connectionGw.GetConnectionRequest{Uuid: id})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	var requests []domain.ConnectionRequest
	for _, user := range connections.Users {
		userResponse, err := userClient.GetOne(context.TODO(), &userGw.GetOneMessage{Id: user.UserID})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
			// ctx.AbortWithError(http.StatusBadRequest, err)
			// return
		}
		requests = append(requests, loadUserInfo(userResponse.User))
	}

	// response := feed
	response, err := json.Marshal(requests)
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

func loadUserInfo(userPb *userGw.User) domain.ConnectionRequest {
	var request domain.ConnectionRequest

	request.UserId = userPb.Id
	request.Name = userPb.Name
	request.Surname = userPb.Surname
	request.Biography = userPb.Biography

	return request
}
