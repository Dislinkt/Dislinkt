package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/dislinkt/common/tracer"
	uuid "github.com/gofrs/uuid"

	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/domain"
)

type ConnectionService struct {
	store domain.ConnectionStore
}

func NewConnectionService(store domain.ConnectionStore) *ConnectionService {
	return &ConnectionService{
		store: store,
	}
}

func (service *ConnectionService) Register(ctx context.Context, userID string, status string) (*domain.UserNode, error) {
	span := tracer.StartSpanFromContext(ctx, "Register-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService Register]")
	if !IsValidUUID(userID) {
		return nil, errors.New("Invalid uuid")
	}
	node := domain.UserNode{UserUID: userID, Status: domain.ProfileStatus(status)}
	return service.store.Register(&node)
}

func (service *ConnectionService) CreateConnection(ctx context.Context, baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateConnection-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService CreateConnection]")
	if !IsValidUUID(baseUserUuid) {
		return nil, errors.New("Invalid uuid")
	}
	return service.store.CreateConnection(baseUserUuid, connectUserUuid)
}

func (service *ConnectionService) AcceptConnection(ctx context.Context, requestSenderUser string, requestApprovalUser string) (*pb.NewConnectionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "AcceptConnection-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService AcceptConnection")
	return service.store.AcceptConnection(requestSenderUser, requestApprovalUser)
}

func (service *ConnectionService) GetAllConnectionForUser(ctx context.Context, userUid string) ([]*domain.UserNode, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConnectionForUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(userUid) {
		return nil, errors.New("Invalid uuid")
	}
	return service.store.GetAllConnectionForUser(userUid)
}

func (service *ConnectionService) GetAllConnectionRequestsForUser(ctx context.Context, userUid string) ([]*domain.UserNode, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConnectionRequestsForUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(userUid) {
		return nil, errors.New("Invalid uuid")
	}
	return service.store.GetAllConnectionRequestsForUser(userUid)
}

func (service *ConnectionService) UpdateUser(ctx context.Context, userUUID string,
	private bool) error {
	span := tracer.StartSpanFromContext(ctx, "UpdateUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(userUUID) {
		return errors.New("Invalid uuid")
	}
	fmt.Println("[ConnectionService UpdateUser")
	return service.store.UpdateUser(userUUID, private)
}

func (service *ConnectionService) BlockUser(ctx context.Context, currentUserUUID string, blockedUserUUid string) (*pb.BlockedUserStatus, error) {
	span := tracer.StartSpanFromContext(ctx, "BlockUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService BlockUser")
	return service.store.BlockUser(currentUserUUID, blockedUserUUid)
}

func (servioe *ConnectionService) GetAllBlockedForCurrentUser(ctx context.Context, currentUserUUID string) ([]*domain.UserNode, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllBlockedForCurrentUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService GetAllBlockedForCurrentUser")
	return servioe.store.GetAllBlockedForCurrentUser(currentUserUUID)
}

func (servioe *ConnectionService) GetAllUserBlockingCurrentUser(ctx context.Context, currentUserUUID string) ([]*domain.UserNode, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllUserBlockingCurrentUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService GetAllUserBlockingCurrentUser")
	return servioe.store.GetAllUserBlockingCurrentUser(currentUserUUID)
}

func (servioe *ConnectionService) RecommendUsersByConnection(ctx context.Context, currentUserUUID string) (users []*domain.UserNode, err error) {
	span := tracer.StartSpanFromContext(ctx, "RecommendUsersByConnection-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService RecommendUsersByConnection")
	return servioe.store.RecommendUsersByConnection(currentUserUUID)
}

func (servioe *ConnectionService) UnblockConnection(ctx context.Context, currentUser string, blockedUser string) (*pb.BlockedUserStatus, error) {
	span := tracer.StartSpanFromContext(ctx, "UnblockConnection-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService UnblockConnection")
	return servioe.store.UnblockConnection(currentUser, blockedUser)
}

func (servioe *ConnectionService) InsertField(ctx context.Context, name string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "InsertField-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService InsertField")
	return servioe.store.InsertField(name)
}

func (servioe *ConnectionService) InsertSkill(ctx context.Context, name string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "InsertSkill-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService InsertSkill")
	return servioe.store.InsertSkill(name)
}

func (servioe *ConnectionService) InsertJobOffer(ctx context.Context, jobOffer domain.JobOffer) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "InsertJobOffer-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService InsertJobOffer")
	return servioe.store.InsertJobOffer(jobOffer)
}

func (servioe *ConnectionService) InsertSkillToUser(ctx context.Context, name string, uuid string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "InsertSkillToUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService InsertSkillToUser")
	fmt.Println(name)
	fmt.Println(uuid)
	return servioe.store.InsertSkillToUser(name, uuid)
}

func (servioe *ConnectionService) UpdateSkillForUser(ctx context.Context, userUUID string, skillNameOld string, skillNameNew string) (res string, err error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserSkill-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService UpdateUserSkill")

	return servioe.store.UpdateSkillForUser(userUUID, skillNameOld, skillNameNew)
}

func (servioe *ConnectionService) UpdateEducationForUser(ctx context.Context, userUUID string, educationNameOld string, educationNameNew string) (res string, err error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateEducationForUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService UpdateEducationForUser")

	return servioe.store.UpdateFieldForUser(userUUID, educationNameOld, educationNameNew)
}

func (servioe *ConnectionService) DeleteSkillToUser(ctx context.Context, name string, uuid string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteSkillToUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService DeleteSkillToUser")
	fmt.Println(name)
	fmt.Println(uuid)
	return servioe.store.DeleteSkillForUser(uuid, name)
}

func (servioe *ConnectionService) DeleteFieldToUser(ctx context.Context, name string, uuid string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteFieldToUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService DeleteFieldToUser")
	fmt.Println(name)
	fmt.Println(uuid)
	return servioe.store.DeleteFieldForUser(uuid, name)
}

func (servioe *ConnectionService) InsertFieldToUser(ctx context.Context, name string, uuid string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "InsertFieldToUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService InsertFieldToUser")
	fmt.Println(name, uuid)
	return servioe.store.InsertFieldToUser(name, uuid)
}

func (servioe *ConnectionService) RecommendJobBySkill(ctx context.Context, userUid string) (jobs []*domain.JobOffer, err error) {
	span := tracer.StartSpanFromContext(ctx, "RecommendJobBySkill-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService RecommendJobBySkill")
	return servioe.store.RecommendJobBySkill(userUid)
}

func (servioe *ConnectionService) RecommendJobByField(ctx context.Context, userUid string) (jobs []*domain.JobOffer, err error) {
	span := tracer.StartSpanFromContext(ctx, "RecommendJobByField-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService RecommendJobByField11")
	return servioe.store.RecommendJobByField(userUid)
}

func (servioe *ConnectionService) CheckIfUsersConnected(ctx context.Context, userUUID1 string, userUUID2 string) (bool, error) {
	span := tracer.StartSpanFromContext(ctx, "CheckIfUsersConnected-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService CHeckIfUsersConnected")
	return servioe.store.CheckIfUsersConnected(userUUID1, userUUID2)
}

func (servioe *ConnectionService) CheckIfUsersBlocked(ctx context.Context, userUUID1 string, userUUID2 string) (bool, error) {
	span := tracer.StartSpanFromContext(ctx, "CheckIfUserBlocked-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	fmt.Println("[ConnectionService CheckIfUsersBlocked")
	return servioe.store.CheckIfUsersBlocked(userUUID1, userUUID2)
}

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}
