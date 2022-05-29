package config

func AccessibleRoles() map[string][]string {
	const userService = "/user_service_proto.UserService/"

	return map[string][]string{
		userService + "UpdateUser": {"Regular"},
		userService + "PatchUser":  {"Regular"},
	}
}
