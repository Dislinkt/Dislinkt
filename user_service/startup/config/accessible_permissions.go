package config

func AccessiblePermissions() map[string][]string {
	const userService = "/user_service_proto.UserService/"

	return map[string][]string{
		userService + "UpdateUser": {"updateUserPermission"},
		userService + "PatchUser":  {"patchUserPermission"},
	}
}
