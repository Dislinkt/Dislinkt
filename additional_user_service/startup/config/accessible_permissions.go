package config

func AccessiblePermissions() map[string][]string {
	const additionalUserService = "/additional_user_service_proto.AdditionalUserService/"

	return map[string][]string{
		additionalUserService + "NewEducation":    {"addNewEducationPermission"},
		additionalUserService + "GetAllEducation": {"getAllEducationPermission"},
		additionalUserService + "UpdateEducation": {"updateEducationPermission"},
		additionalUserService + "DeleteEducation": {"deleteEducationPermission"},
		additionalUserService + "NewPosition":     {"addNewPositionPermission"},
		additionalUserService + "GetAllPosition":  {"getAllPositionPermission"},
		additionalUserService + "UpdatePosition":  {"updatePositionPermission"},
		additionalUserService + "DeletePosition":  {"deletePositionPermission"},
		additionalUserService + "NewSkill":        {"addNewSkillPermission"},
		additionalUserService + "GetAllSkill":     {"getAllSkillPermission"},
		additionalUserService + "UpdateSkill":     {"updateSkillPermission"},
		additionalUserService + "DeleteSkill":     {"deleteSkillPermission"},
		additionalUserService + "NewInterest":     {"addNewInterestPermission"},
		additionalUserService + "GetAllInterest":  {"getAllInterestPermission"},
		additionalUserService + "UpdateInterest":  {"updateInterestPermission"},
		additionalUserService + "DeleteInterest":  {"deleteInterestPermission"},
	}
}
