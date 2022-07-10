package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/interceptor"
	connectionGw "github.com/dislinkt/common/proto/connection_service"
	pb "github.com/dislinkt/common/proto/notification_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/common/tracer"
	"github.com/dislinkt/notification_service/application"
	"github.com/dislinkt/notification_service/domain"
	"github.com/dislinkt/notification_service/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	service *application.NotificationService
}

func NewNotificationHandler(service *application.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (handler *NotificationHandler) GetNotificationsForUser(ctx context.Context, request *pb.Empty) (*pb.GetMultipleResponse, error) {
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	span := tracer.StartSpanFromContext(ctx, "GetNotificationsForUserAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	userResponse, _ := persistence.UserClient("user_service:8000").GetUserByUsername(context.TODO(), &userGw.GetOneByUsernameMessage{Username: username})
	notifications, err := handler.service.GetNotificationsForUser(ctx, userResponse.User.Id)
	if err != nil {
		return nil, err
	}
	response := &pb.GetMultipleResponse{Notifications: []*pb.Notification{}, UnreadNotificationsNumber: 0}
	unreadNotificationsNumber := 0
	for _, notification := range notifications {
		current := mapNotification(notification)
		response.Notifications = append(response.Notifications, current)
		if !notification.IsRead {
			unreadNotificationsNumber++
		}
	}
	response.UnreadNotificationsNumber = int32(unreadNotificationsNumber)
	return response, nil
}

func (handler *NotificationHandler) SaveNotification(ctx context.Context, request *pb.SaveNotificationRequest) (*pb.Empty, error) {
	span := tracer.StartSpanFromContext(ctx, "SaveNotificationAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	notification := mapNewNotification(request.Notification)
	if notification.NotificationType == domain.POST {
		handler.saveNotificationsForNewPost(ctx, notification, request.UserId, request.Notification.SubjectUsername)
	} else {
		userNotificationSettings, _ := persistence.UserClient("user_service:8000").GetNotificationSettings(ctx, &userGw.GetOneMessage{Id: request.UserId})
		if (notification.NotificationType == domain.CONNECTION && userNotificationSettings.ConnectionNotifications) || (notification.NotificationType == domain.MESSAGE && userNotificationSettings.MessageNotifications) {
			notification.UserId = request.UserId
			notification.NotificationText = generateNotificationText(notification, request.Notification.SubjectUsername)
			err := handler.service.InsertNotification(ctx, notification)
			if err != nil {
				return nil, err
			}
		}
	}
	return &pb.Empty{}, nil
}

func (handler *NotificationHandler) saveNotificationsForNewPost(ctx context.Context, notification *domain.Notification, userId string, username string) {
	span := tracer.StartSpanFromContext(ctx, "SaveNotificationForNewPostAPI")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	connections, _ := persistence.ConnectionClient("connection_service:8000").GetAllConnectionForUser(ctx, &connectionGw.GetConnectionRequest{Uuid: userId})
	notification.NotificationText = generateNotificationText(notification, username)

	for _, user := range connections.Users {
		userNotificationSettings, _ := persistence.UserClient("user_service:8000").GetNotificationSettings(ctx, &userGw.GetOneMessage{Id: user.UserID})
		if notification.NotificationType == domain.POST && userNotificationSettings.PostNotifications {
			notification.Id = primitive.NewObjectID()
			notification.UserId = user.UserID
			_ = handler.service.InsertNotification(ctx, notification)
		}
	}
}

func generateNotificationText(notification *domain.Notification, subjectUsername string) string {
	if notification.NotificationType == domain.CONNECTION {
		return subjectUsername + " added you as a connection."
	} else if notification.NotificationType == domain.MESSAGE {
		return subjectUsername + " sent you a message."
	} else if notification.NotificationType == domain.POST {
		return subjectUsername + " made a new post."
	} else {
		return subjectUsername + " sent you a connection request."
	}
}
