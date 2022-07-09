package api

import (
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/common/saga/events"
)

func mapAdditionalCommandAddSkill(command *events.AddSkillCommand) *domain.Skill {

	skillD := &domain.Skill{
		Id:   command.Skill.Id,
		Name: command.Skill.Name,
	}
	return skillD
}
