package config

func AccessibleRoles() map[string][]string {
	const authService = "/proto.AuthService/"

	return map[string][]string{
		authService + "ChangePassword": {"Regular"},
	}
}
