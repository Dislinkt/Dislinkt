package api

import (
	"context"
	"fmt"

	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/application"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	service *application.ConnectionService
}

func NewConnectionHandler(service *application.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		service: service,
	}
}

func (handler *ConnectionHandler) Register(ctx context.Context, request *pb.RegisterRequest) (response *pb.RegisterResponse, err error) {
	fmt.Println("[ConnectionHandler]:Register")
	userID := request.User.UserID
	status := request.User.Status
	item, err := handler.service.Register(userID, status)
	message := ""

	if item.Status == "" {
		message = "User already exists"
	} else {
		message = "Success"
	}

	user := &pb.User{
		UserID: userID, Status: status,
	}
	actionResult := &pb.RegisterResponse{User: user, Message: message}

	return actionResult, err
}

func (handler *ConnectionHandler) CreateConnection(ctx context.Context, request *pb.NewConnectionRequest) (response *pb.NewConnectionResponse, err error) {
	fmt.Println("[ConnectionHandler]:CreateConnection")
	return handler.service.CreateConnection(request.Connection.BaseUserUUID, request.Connection.ConnectUserUUID)
}

func (handler *ConnectionHandler) AcceptConnection(ctx context.Context, request *pb.AcceptConnectionMessage) (response *pb.NewConnectionResponse, err error) {
	fmt.Println("[ConnectionHandler]:AcceptConnection")
	return handler.service.AcceptConnection(request.AcceptConnection.RequestSenderUser, request.AcceptConnection.RequestApprovalUser)
}

func (handler *ConnectionHandler) GetAllConnectionForUser(ctx context.Context, request *pb.GetConnectionRequest) (*pb.GetAllResponse, error) {

	users, err := handler.service.GetAllConnectionForUser(request.GetUuid())
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := pb.User{UserID: user.UserUID, Status: string(user.Status)}
		response.Users = append(response.Users, &current)
	}

	return response, nil
}

func (handler *ConnectionHandler) GetAllConnectionRequestsForUser(ctx context.Context,
	request *pb.GetConnectionRequest) (*pb.GetAllResponse, error) {

	users, err := handler.service.GetAllConnectionRequestsForUser(request.GetUuid())
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := pb.User{UserID: user.UserUID, Status: string(user.Status)}
		response.Users = append(response.Users, &current)
	}

	return response, nil
}

func (handler *ConnectionHandler) BlockUser(ctx context.Context, request *pb.BlockUserRequest) (response *pb.BlockedUserStatus, err error) {
	fmt.Println("[ConnectionHandler]:BlockUser")
	return handler.service.BlockUser(request.Uuid, request.Uuid1)
}

func (handler *ConnectionHandler) GetAllBlockedForCurrentUser(ctx context.Context, request *pb.BlockUserRequest) (*pb.GetAllResponse, error) {
	fmt.Println("[ConnectionHandler]:GetAllBlockedForCurrentUser")
	users, err := handler.service.GetAllBlockedForCurrentUser(request.Uuid)
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		fmt.Println("Uslo je")
		current := pb.User{UserID: user.UserUID, Status: string(user.Status)}
		response.Users = append(response.Users, &current)
	}

	return response, nil
}

func (handler *ConnectionHandler) GetAllUserBlockingCurrentUser(ctx context.Context, request *pb.BlockUserRequest) (*pb.GetAllResponse, error) {
	fmt.Println("[ConnectionHandler]:GetAllUserBlockingCurrentUser")
	users, err := handler.service.GetAllUserBlockingCurrentUser(request.Uuid)
	fmt.Println(len(users))
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := pb.User{UserID: user.UserUID, Status: string(user.Status)}
		response.Users = append(response.Users, &current)
	}

	return response, nil
}

func (handler *ConnectionHandler) RecommendUsersByConnection(ctx context.Context, request *pb.GetConnectionRequest) (*pb.GetAllResponse, error) {
	fmt.Println("[ConnectionHandler]:RecommendUsersByConnection")
	users, err := handler.service.RecommendUsersByConnection(request.Uuid)
	fmt.Println(len(users))
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := pb.User{UserID: user.UserUID, Status: string(user.Status)}
		response.Users = append(response.Users, &current)
	}

	return response, nil
}

func (handler *ConnectionHandler) UnblockConnection(ctx context.Context, request *pb.BlockUserRequest) (response *pb.BlockedUserStatus, err error) {
	fmt.Println("[ConnectionHandler]:UnblockConnection")
	return handler.service.UnblockConnection(request.Uuid, request.Uuid1)
}

func (handler *ConnectionHandler) InsertField(ctx context.Context, request *pb.Field) (response *pb.Response, err error) {
	fmt.Println("[ConnectionHandler]:InsertField")
	name, err := handler.service.InsertField(request.Name)
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) InsertSkill(ctx context.Context, request *pb.Skill) (response *pb.Response, err error) {
	fmt.Println("[ConnectionHandler]:InsertSkill")
	name, err := handler.service.InsertSkill(request.Name)
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) InsertJobOffer(ctx context.Context, request *pb.JobOffer) (response *pb.Response, err error) {
	fmt.Println("[ConnectionHandler]:InsertJobOffer")
	name, err := handler.service.InsertJobOffer(*mapJobOfferPb(request))
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) InsertSkillToUser(ctx context.Context, request *pb.UserInfoItem) (response *pb.Response, err error) {
	fmt.Println("[ConnectionHandler]:InsertSkillToUser")
	name, err := handler.service.InsertSkillToUser(request.Name, request.Uuid)
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) InsertFieldToUser(ctx context.Context, request *pb.UserInfoItem) (response *pb.Response, err error) {
	fmt.Println("[ConnectionHandler]:InsertFieldToUser")
	name, err := handler.service.InsertFieldToUser(request.Name, request.Uuid)
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) RecommendJobBySkill(ctx context.Context, request *pb.GetConnectionRequest) (response *pb.JobOffers, err error) {
	fmt.Println("[ConnectionHandler]:InsertFieldToUser")
	response = &pb.JobOffers{
		Jobs: []*pb.JobOffer{},
	}
	jobs, err := handler.service.RecommendJobBySkill(request.Uuid)

	for _, job := range jobs {
		fmt.Println("Uslo je")
		current := mapJobOffer(job)
		response.Jobs = append(response.Jobs, current)
	}
	return response, err
}

func (handler *ConnectionHandler) RecommendJobByField(ctx context.Context, request *pb.GetConnectionRequest) (response *pb.JobOffers, err error) {
	fmt.Println("[ConnectionHandler]:RecommendJobByField")
	response = &pb.JobOffers{
		Jobs: []*pb.JobOffer{},
	}
	jobs, err := handler.service.RecommendJobByField(request.Uuid)

	for _, job := range jobs {
		current := mapJobOffer(job)
		response.Jobs = append(response.Jobs, current)
	}
	return response, err
}

func (handler *ConnectionHandler) CheckIfUsersConnected(ctx context.Context, request *pb.CheckConnection) (response *pb.CheckResult, err error) {
	fmt.Println("[ConnectionHandler]:CheckIfUsersConnected")
	isConnected, err := handler.service.CheckIfUsersConnected(request.Uuid1, request.Uuid2)
	if err != nil {
		return nil, err
	}
	response = &pb.CheckResult{
		IsConnected: isConnected,
	}

	return response, err
}
