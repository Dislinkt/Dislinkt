package application

import (
	"github.com/dislinkt/notification_service/domain"
)

type NotificationService struct {
	store domain.NotificationStore
}

func NewNotificationService(store domain.NotificationStore) *NotificationService {
	return &NotificationService{store: store}
}

func (service *NotificationService) GetNotificationsForUser(userId string) ([]*domain.Notification, error) {
	return service.store.GetNotificationsForUser(userId)
}

func (service *NotificationService) InsertNotification(notification *domain.Notification) error {
	return service.store.InsertNotification(notification)
}
