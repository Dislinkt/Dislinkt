package config

func AccessiblePermissions() map[string]string {
	const messageService = "/message_service_proto.MessageService/"

	return map[string]string{
		messageService + "GetMessageHistoriesByUser": "sendMessage",
		messageService + "GetMessageHistory":         "sendMessage",
		messageService + "SendMessage":               "sendMessage",
	}
}
