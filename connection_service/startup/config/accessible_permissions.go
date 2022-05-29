package config

func AccessiblePermissions() map[string]string {
	const connectionService = "/connection_service_proto.ConnectionService/"

	return map[string]string{
		connectionService + "CreateConnection": "createConnectionPermission",
		connectionService + "AcceptConnection": "acceptConnectionPermission",
	}
}
