package application

import (
	"github.com/dislinkt/additional_user_service/domain"
)

type AdditionalUserService struct {
	store domain.AdditionalUserStore
}

func (service *AdditionalUserService) CreateDocument(uuid string) error {
	_, err := service.store.CreateUserDocument(uuid)
	if err != nil {
		return err
	}

	return nil
}

func (service *AdditionalUserService) DeleteDocument(uuid string) error {
	err := service.store.DeleteUserDocument(uuid)
	if err != nil {
		return err
	}

	return nil
}

// EDUCATION

func NewAdditionalUserService(store domain.AdditionalUserStore) *AdditionalUserService {
	return &AdditionalUserService{
		store: store,
	}
}

func (service *AdditionalUserService) CreateEducation(uuid string, education *domain.Education) (*domain.Education,
	error) {

	_, err := service.store.FindDocument(uuid)
	if err != nil {
		return nil, err
	}

	insertEducation, err := service.store.InsertEducation(uuid, education)
	if err != nil {
		return nil, err
	}

	return insertEducation, nil
}

func (service *AdditionalUserService) FindUserEducations(uuid string) (*map[string]domain.Education,
	error) {
	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Educations, nil
}

func (service *AdditionalUserService) UpdateUserEducation(uuid string, educationId string,
	education *domain.Education) (*map[string]domain.Education, error) {

	err := service.store.UpdateUserEducation(educationId, education)
	if err != nil {
		return nil, err
	}
	userEducation, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userEducation.Educations, nil
}

func (service *AdditionalUserService) DeleteUserEducation(uuid string, additionID string) (*map[string]domain.Education, error) {

	err := service.store.DeleteUserEducation(additionID)
	userEducation, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userEducation.Educations, nil
}

// POSITION

func (service *AdditionalUserService) CreatePosition(uuid string, position *domain.Position) (*domain.Position,
	error) {

	_, err := service.store.FindDocument(uuid)
	if err != nil {
		return nil, err
	}

	insertPosition, err := service.store.InsertPosition(uuid, position)
	if err != nil {
		return nil, err
	}

	return insertPosition, nil
}

func (service *AdditionalUserService) FindUserPositions(uuid string) (*map[string]domain.Position, error) {
	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Positions, nil
}

func (service *AdditionalUserService) UpdateUserPosition(uuid string, positionId string,
	position *domain.Position) (*map[string]domain.Position, error) {
	err := service.store.UpdateUserPosition(positionId, position)
	if err != nil {
		return nil, err
	}
	userPosition, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userPosition.Positions, nil
}

func (service *AdditionalUserService) DeleteUserPosition(uuid string, additionID string) (*map[string]domain.Position,
	error) {

	err := service.store.DeleteUserPosition(additionID)
	userPosition, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userPosition.Positions, nil
}

// SKILL

func (service *AdditionalUserService) CreateSkill(uuid string, skill *domain.Skill) (*domain.Skill,
	error) {

	_, err := service.store.FindDocument(uuid)
	if err != nil {
		return nil, err
	}

	insertSkill, err := service.store.InsertSkill(uuid, skill)
	if err != nil {
		return nil, err
	}

	return insertSkill, nil
}

func (service *AdditionalUserService) FindUserSkills(uuid string) (*map[string]domain.Skill, error) {
	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Skills, nil
}

func (service *AdditionalUserService) UpdateUserSkill(uuid string, skillId string,
	skill *domain.Skill) (*map[string]domain.Skill, error) {
	err := service.store.UpdateUserSkill(skillId, skill)
	if err != nil {
		return nil, err
	}
	userSkill, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userSkill.Skills, nil
}

func (service *AdditionalUserService) DeleteUserSkill(uuid string, additionID string) (*map[string]domain.Skill,
	error) {

	err := service.store.DeleteUserSkill(additionID)
	userSkill, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userSkill.Skills, nil
}

// INTEREST

func (service *AdditionalUserService) CreateInterest(uuid string, interest *domain.Interest) (*domain.Interest,
	error) {

	_, err := service.store.FindDocument(uuid)
	if err != nil {
		return nil, err
	}

	insertInterest, err := service.store.InsertInterest(uuid, interest)
	if err != nil {
		return nil, err
	}

	return insertInterest, nil
}

func (service *AdditionalUserService) FindUserInterests(uuid string) (*map[string]domain.Interest, error) {
	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Interests, nil
}

func (service *AdditionalUserService) UpdateUserInterest(uuid string, interestId string,
	interest *domain.Interest) (*map[string]domain.Interest, error) {
	err := service.store.UpdateUserInterest(interestId, interest)
	if err != nil {
		return nil, err
	}
	userInterest, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userInterest.Interests, nil
}

func (service *AdditionalUserService) DeleteUserInterest(uuid string, additionID string) (*map[string]domain.Interest,
	error) {

	err := service.store.DeleteUserInterest(additionID)
	userInterest, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userInterest.Interests, nil
}
