package config

func AccessibleRoles() map[string][]string {
	const messageService = "/message_service_proto.MessageService/"

	return map[string][]string{
		messageService + "GetMessageHistoriesByUser": {"Regular"},
		messageService + "GetMessageHistory":         {"Regular"},
		messageService + "SendMessage":               {"Regular"},
	}
}
