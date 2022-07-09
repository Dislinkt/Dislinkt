package api

import (
	pb "github.com/dislinkt/common/proto/notification_service"
	"github.com/dislinkt/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func mapNotification(notification *domain.Notification) *pb.Notification {
	notificationPb := &pb.Notification{
		NotificationText: notification.NotificationText,
		NotificationType: pb.NotificationType(notification.NotificationType),
		Date:             notification.Date.String(),
		IsRead:           notification.IsRead,
	}

	if notificationPb.NotificationType == 0 {
		notificationPb.NotificationType = 1
	}

	return notificationPb
}

func mapNewNotification(notificationPb *pb.NewNotification) *domain.Notification {
	notification := &domain.Notification{
		Id:               primitive.NewObjectID(),
		NotificationType: domain.NotificationType(notificationPb.NotificationType),
		Date:             time.Now(),
		IsRead:           false,
	}

	return notification
}
