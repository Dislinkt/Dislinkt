package api

import (
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/common/saga/events"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mapAdditionalCommandAddSkill(command *events.AddSkillCommand) *domain.Skill {

	skillD := &domain.Skill{
		Id:   command.Skill.Id,
		Name: command.Skill.Name,
	}
	return skillD
}

func mapAdditionalCommandUpdateSkill(command *events.UpdateSkillCommand) *domain.Skill {

	skillD := &domain.Skill{
		Id:   primitive.NewObjectID(),
		Name: command.Skill.Name,
	}
	return skillD
}
