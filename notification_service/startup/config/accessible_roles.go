package config

func AccessibleRoles() map[string][]string {
	const notificationService = "/notification_service_proto.NotificationService/"

	return map[string][]string{
		notificationService + "GetNotificationsForUser": {"Regular"},
	}
}
