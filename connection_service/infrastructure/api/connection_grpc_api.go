package api

import (
	"context"
	"fmt"
	eventGw "github.com/dislinkt/common/proto/event_service"
	notificationGw "github.com/dislinkt/common/proto/notification_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/common/tracer"
	"github.com/dislinkt/connection_service/infrastructure/persistance"

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
	span := tracer.StartSpanFromContext(ctx, "RegisterAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:Register")
	userID := request.User.UserID
	status := request.User.Status
	item, err := handler.service.Register(ctx, userID, status)
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
	span := tracer.StartSpanFromContext(ctx, "CreateConnectionAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:CreateConnection")
	response, err = handler.service.CreateConnection(ctx, request.Connection.BaseUserUUID, request.Connection.ConnectUserUUID)

	userResponse, _ := persistance.UserClient("user_service:8000").GetOne(context.TODO(), &userGw.GetOneMessage{Id: request.Connection.ConnectUserUUID})
	if response.ConnectionResponse == "CONNECTED" {
		_, _ = persistance.EventClient("event_service:8000").SaveEvent(context.TODO(),
			&eventGw.SaveEventRequest{Event: mapEventForConnection(request.Connection.ConnectUserUUID, request.Connection.BaseUserUUID)})
	} else {
		_, _ = persistance.EventClient("event_service:8000").SaveEvent(context.TODO(),
			&eventGw.SaveEventRequest{Event: mapEventForConnectionRequest(request.Connection.ConnectUserUUID, request.Connection.BaseUserUUID)})
	}
	_, _ = persistance.NotificationClient("notification_service:8000").SaveNotification(context.TODO(),
		&notificationGw.SaveNotificationRequest{Notification: mapNotification(userResponse.User.Username, response.ConnectionResponse), UserId: request.Connection.BaseUserUUID})

	return response, err
}

func (handler *ConnectionHandler) AcceptConnection(ctx context.Context, request *pb.AcceptConnectionMessage) (response *pb.NewConnectionResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "AcceptConnestionAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:AcceptConnection")
	return handler.service.AcceptConnection(ctx, request.AcceptConnection.RequestSenderUser, request.AcceptConnection.RequestApprovalUser)
}

func (handler *ConnectionHandler) GetAllConnectionForUser(ctx context.Context, request *pb.GetConnectionRequest) (*pb.GetAllResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConnectionForUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	users, err := handler.service.GetAllConnectionForUser(ctx, request.GetUuid())
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
	span := tracer.StartSpanFromContext(ctx, "GetAllConnectionRequestsForUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	users, err := handler.service.GetAllConnectionRequestsForUser(ctx, request.GetUuid())
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
	span := tracer.StartSpanFromContext(ctx, "BlockUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:BlockUser")
	_, _ = persistance.EventClient("event_service:8000").SaveEvent(context.TODO(),
		&eventGw.SaveEventRequest{Event: mapEventForUserBlocking(request.Uuid, request.Uuid1)})
	return handler.service.BlockUser(ctx, request.Uuid, request.Uuid1)
}

func (handler *ConnectionHandler) GetAllBlockedForCurrentUser(ctx context.Context, request *pb.BlockUserRequest) (*pb.GetAllResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllBlockedForCurrentUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:GetAllBlockedForCurrentUser")
	users, err := handler.service.GetAllBlockedForCurrentUser(ctx, request.Uuid)
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
	span := tracer.StartSpanFromContext(ctx, "GetAllUserBlockingCurrentUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:GetAllUserBlockingCurrentUser")
	users, err := handler.service.GetAllUserBlockingCurrentUser(ctx, request.Uuid)
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
	span := tracer.StartSpanFromContext(ctx, "RecommendUsersByConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:RecommendUsersByConnection")
	users, err := handler.service.RecommendUsersByConnection(ctx, request.Uuid)
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
	span := tracer.StartSpanFromContext(ctx, "UnblockConnectionaPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:UnblockConnection")
	return handler.service.UnblockConnection(ctx, request.Uuid, request.Uuid1)
}

func (handler *ConnectionHandler) InsertField(ctx context.Context, request *pb.Field) (response *pb.Response, err error) {
	span := tracer.StartSpanFromContext(ctx, "InsertFieldAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:InsertField")
	name, err := handler.service.InsertField(ctx, request.Name)
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) InsertSkill(ctx context.Context, request *pb.Skill) (response *pb.Response, err error) {
	span := tracer.StartSpanFromContext(ctx, "InsertSkillAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:InsertSkill")
	name, err := handler.service.InsertSkill(ctx, request.Name)
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) InsertJobOffer(ctx context.Context, request *pb.JobOffer) (response *pb.Response, err error) {
	span := tracer.StartSpanFromContext(ctx, "InsertJobOfferAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:InsertJobOffer")
	name, err := handler.service.InsertJobOffer(ctx, *mapJobOfferPb(request))
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) InsertSkillToUser(ctx context.Context, request *pb.UserInfoItem) (response *pb.Response, err error) {
	span := tracer.StartSpanFromContext(ctx, "InsertSkillToUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:InsertSkillToUser")
	name, err := handler.service.InsertSkillToUser(ctx, request.Name, request.Uuid)
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) InsertFieldToUser(ctx context.Context, request *pb.UserInfoItem) (response *pb.Response, err error) {
	span := tracer.StartSpanFromContext(ctx, "InsertFieldToUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:InsertFieldToUser")
	name, err := handler.service.InsertFieldToUser(ctx, request.Name, request.Uuid)
	response = &pb.Response{
		Name: name,
	}
	return response, err
}

func (handler *ConnectionHandler) RecommendJobBySkill(ctx context.Context, request *pb.GetConnectionRequest) (response *pb.JobOffers, err error) {
	span := tracer.StartSpanFromContext(ctx, "RecommendJobBySkill")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:RecommendJobBySkill")
	response = &pb.JobOffers{
		Jobs: []*pb.JobOffer{},
	}
	jobs, err := handler.service.RecommendJobBySkill(ctx, request.Uuid)

	for _, job := range jobs {
		fmt.Println("Uslo je")
		current := mapJobOffer(job)
		response.Jobs = append(response.Jobs, current)
	}
	return response, err
}

func (handler *ConnectionHandler) RecommendJobByField(ctx context.Context, request *pb.GetConnectionRequest) (response *pb.JobOffers, err error) {
	span := tracer.StartSpanFromContext(ctx, "RecommendJobByField")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:RecommendJobByField")
	response = &pb.JobOffers{
		Jobs: []*pb.JobOffer{},
	}
	jobs, err := handler.service.RecommendJobByField(ctx, request.Uuid)

	for _, job := range jobs {
		current := mapJobOffer(job)
		response.Jobs = append(response.Jobs, current)
	}
	return response, err
}

func (handler *ConnectionHandler) CheckIfUsersConnected(ctx context.Context, request *pb.CheckConnection) (response *pb.CheckResult, err error) {
	span := tracer.StartSpanFromContext(ctx, "CheckIfUsersConnected")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:CheckIfUsersConnected")
	isConnected, err := handler.service.CheckIfUsersConnected(ctx, request.Uuid1, request.Uuid2)
	if err != nil {
		return nil, err
	}
	response = &pb.CheckResult{
		IsConnected: isConnected,
	}

	return response, err
}

func (handler *ConnectionHandler) CheckIfUsersBlocked(ctx context.Context, request *pb.CheckConnection) (response *pb.CheckResultBlock, err error) {
	span := tracer.StartSpanFromContext(ctx, "CheckIfUsersBlocked")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionHandler]:CheckIfUsersConnected")
	isBlocked, err := handler.service.CheckIfUsersBlocked(ctx, request.Uuid1, request.Uuid2)
	if err != nil {
		return nil, err
	}
	response = &pb.CheckResultBlock{
		IsBlocked: isBlocked,
	}

	return response, err
}
