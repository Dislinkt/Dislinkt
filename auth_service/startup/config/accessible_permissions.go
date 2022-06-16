package config

func AccessiblePermissions() map[string]string {
	const authService = "/proto.AuthService/"

	return map[string]string{
		authService + "ChangePassword": "changePasswordPermission",
		authService + "Get2FA":         "get2faPermission",
		authService + "Set2FA":         "set2faPermission",
	}
}
