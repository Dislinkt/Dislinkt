package domain

type NotificationStore interface {
	GetNotificationsForUser(userId string) ([]*Notification, error)
	InsertNotification(notification *Notification) error
}
