package config

func AccessiblePermissions() map[string]string {
	const eventService = "/event_service_proto.EventService/"

	return map[string]string{
		eventService + "GetAllEvents": "readEvents",
	}
}
