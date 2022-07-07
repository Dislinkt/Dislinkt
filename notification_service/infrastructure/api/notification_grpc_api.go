package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/interceptor"
	pb "github.com/dislinkt/common/proto/notification_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"github.com/dislinkt/notification_service/application"
	"github.com/dislinkt/notification_service/domain"
	"github.com/dislinkt/notification_service/infrastructure/persistence"
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
	userResponse, _ := persistence.UserClient("user_service:8000").GetUserByUsername(context.TODO(), &userGw.GetOneByUsernameMessage{Username: username})
	notifications, err := handler.service.GetNotificationsForUser(userResponse.User.Id)
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
	notification := mapNewNotification(request.Notification)
	notification.UserId = request.UserId
	notification.NotificationText = generateNotificationText(notification, request.Notification.SubjectUsername)
	err := handler.service.InsertNotification(notification)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func generateNotificationText(notification *domain.Notification, subjectUsername string) string {
	if notification.NotificationType == domain.CONNECTION {
		return subjectUsername + " added you as a connection."
	} else if notification.NotificationType == domain.MESSAGE {
		return subjectUsername + " sent you a message."
	} else if notification.NotificationType == domain.POST {
		return subjectUsername + " made a new post."
	}
	return ""
}
