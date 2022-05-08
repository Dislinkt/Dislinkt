package routes

import (
	"context"
	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AcceptConnectionRequestBody struct {
	RequestSenderUser   string `json:"request_sender_user"`
	RequestApprovalUser string `json:"request_approval_user"`
}

func AcceptConnection(ctx *gin.Context, c pb.ConnectionServiceClient) {
	b := AcceptConnectionRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.AcceptConnection(context.Background(), &pb.AcceptConnectionMessage{
		AcceptConnection: &pb.AcceptConnection{
			RequestSenderUser:   b.RequestSenderUser,
			RequestApprovalUser: b.RequestApprovalUser},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
