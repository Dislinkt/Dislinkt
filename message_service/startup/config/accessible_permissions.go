package config

func AccessiblePermissions() map[string]string {
	const messageService = "/message_service_proto.MessageService/"

	return map[string]string{
		messageService + "getMessageHistoriesByUser": "sendMessage",
		messageService + "getMessageHistory":         "sendMessage",
		messageService + "sendMessage":               "sendMessage",
	}
}
