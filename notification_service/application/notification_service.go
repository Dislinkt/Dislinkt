package application

import (
	"context"
	"github.com/dislinkt/common/tracer"
	"github.com/dislinkt/notification_service/domain"
)

type NotificationService struct {
	store domain.NotificationStore
}

func NewNotificationService(store domain.NotificationStore) *NotificationService {
	return &NotificationService{store: store}
}

func (service *NotificationService) GetNotificationsForUser(ctx context.Context, userId string) ([]*domain.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "GetNotificationsForUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetNotificationsForUser(userId)
}

func (service *NotificationService) InsertNotification(ctx context.Context, notification *domain.Notification) error {
	span := tracer.StartSpanFromContext(ctx, "InsertNotification-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.InsertNotification(notification)
}
