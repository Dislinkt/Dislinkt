package config

func AccessiblePermissions() map[string]string {
	const authService = "/proto.AuthService/"

	return map[string]string{
		authService + "ChangePassword": "changePasswordPermission",
	}
}
