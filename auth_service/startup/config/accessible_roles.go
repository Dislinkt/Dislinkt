package config

// => /naziv_paketa.naziv_servisa/naziv_metode : { naziv_role1, naziv_role2 }
// za sve metode koje ne treba da se presrecu -> ne dodaju se u mapu
func AccessibleRoles() map[string][]string {
	const authService = "/proto.AuthService/"

	return map[string][]string{
		authService + "ChangePassword": {"Regular"},
	}
}
