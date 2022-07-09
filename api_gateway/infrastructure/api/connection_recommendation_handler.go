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

type ConnectionRecommendationHandler struct {
	userClientAddress       string
	connectionClientAddress string
}

func NewConnectionRecommendationHandler(c *config.Config) *ConnectionRecommendationHandler {
	return &ConnectionRecommendationHandler{
		userClientAddress:       fmt.Sprintf("%s:%s", c.UserHost, c.UserPort),
		connectionClientAddress: fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort),
	}
}

func (handler *ConnectionRecommendationHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/connection/recommend/users/{userId}", handler.GetRecommendation)
	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionRecommendationHandler) GetRecommendation(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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
	connections, err := connectionClient.RecommendUsersByConnection(context.TODO(), &connectionGw.GetConnectionRequest{Uuid: id})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	var requests []domain.ConnectionUser
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
