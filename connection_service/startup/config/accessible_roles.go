package config

func AccessibleRoles() map[string][]string {
	const connectionService = "/connection_service_proto.ConnectionService/"

	return map[string][]string{
		connectionService + "CreateConnection": {"Regular"},
		connectionService + "AcceptConnection": {"Regular"},
	}
}
