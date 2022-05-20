package config

// => /naziv_paketa.naziv_servisa/naziv_metode : { naziv_role1, naziv_role2 }
// za sve metode koje ne treba da se presrecu -> ne dodaju se u mapu
func AccessibleRoles() map[string][]string {
	const additionalUserService = "/additional_user_service_proto.AdditionalUserService/"

	return map[string][]string{
		additionalUserService + "NewEducation":    {"Regular"},
		additionalUserService + "GetAllEducation": {"Regular"},
		additionalUserService + "UpdateEducation": {"Regular"},
		additionalUserService + "DeleteEducation": {"Regular"},
		additionalUserService + "NewPosition":     {"Regular"},
		additionalUserService + "GetAllPosition":  {"Regular"},
		additionalUserService + "UpdatePosition":  {"Regular"},
		additionalUserService + "DeletePosition":  {"Regular"},
		additionalUserService + "NewSkill":        {"Regular"},
		additionalUserService + "GetAllSkill":     {"Regular"},
		additionalUserService + "UpdateSkill":     {"Regular"},
		additionalUserService + "DeleteSkill":     {"Regular"},
		additionalUserService + "NewInterest":     {"Regular"},
		additionalUserService + "GetAllInterest":  {"Regular"},
		additionalUserService + "UpdateInterest":  {"Regular"},
		additionalUserService + "DeleteInterest":  {"Regular"},
	}
}
