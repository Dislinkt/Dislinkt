package config

func AccessiblePermissions() map[string]string {
	const notificationService = "/notification_service_proto.NotificationService/"

	return map[string]string{
		notificationService + "GetNotificationsForUser": "readNotifications",
	}
}
