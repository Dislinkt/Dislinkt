package application

import (
	"errors"
	"fmt"
	"github.com/dislinkt/additional_user_service/domain"
	uuid "github.com/gofrs/uuid"
)

type AdditionalUserService struct {
	store                    domain.AdditionalUserStore
	addEducationOrchestrator *AddEducationOrchestrator
}

func (service *AdditionalUserService) CreateDocument(uuid string) error {
	if !IsValidUUID(uuid) {
		return errors.New("Invalid uuid")
	}
	_, err := service.store.CreateUserDocument(uuid)
	if err != nil {
		return err
	}

	return nil
}

func (service *AdditionalUserService) DeleteDocument(uuid string) error {
	if !IsValidUUID(uuid) {
		return errors.New("Invalid uuid")
	}

	err := service.store.DeleteUserDocument(uuid)
	if err != nil {
		return err
	}

	return nil
}

// EDUCATION

func NewAdditionalUserService(store domain.AdditionalUserStore, addEducationOrchestrator *AddEducationOrchestrator) *AdditionalUserService {
	return &AdditionalUserService{
		store:                    store,
		addEducationOrchestrator: addEducationOrchestrator,
	}
}

func (service *AdditionalUserService) CreateEducation(uuid string, education *domain.Education) (*domain.Education,
	error) {
	if !IsValidUUID(uuid) {
		fmt.Println("invalid uuid")
		return nil, errors.New("Invalid uuid")
	}

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

func (service *AdditionalUserService) CreateEducationStart(uuid string, education *domain.Education) error {
	err := service.addEducationOrchestrator.Start(education, uuid)
	if err != nil {
		return err
	}
	return err
}

func (service *AdditionalUserService) FindUserEducations(uuid string) (*map[string]domain.Education,
	error) {
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Educations, nil
}

func (service *AdditionalUserService) UpdateUserEducation(uuid string, educationId string,
	education *domain.Education) (*map[string]domain.Education, error) {
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Positions, nil
}

func (service *AdditionalUserService) UpdateUserPosition(uuid string, positionId string,
	position *domain.Position) (*map[string]domain.Position, error) {
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Skills, nil
}

func (service *AdditionalUserService) UpdateUserSkill(uuid string, skillId string,
	skill *domain.Skill) (*map[string]domain.Skill, error) {
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Interests, nil
}

func (service *AdditionalUserService) UpdateUserInterest(uuid string, interestId string,
	interest *domain.Interest) (*map[string]domain.Interest, error) {
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

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
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	err := service.store.DeleteUserInterest(additionID)
	userInterest, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userInterest.Interests, nil
}

func (service *AdditionalUserService) GetFieldOfStudies() ([]*domain.FieldOfStudy, error) {
	study, err := service.store.GetAllFieldOfStudy()
	if err != nil {
		return nil, nil
	}
	return study, nil
}

func (service *AdditionalUserService) GetSkills() ([]*domain.Skill, error) {
	skills, err := service.store.GetSkills()
	if err != nil {
		return nil, nil
	}
	return skills, nil
}

func (service *AdditionalUserService) GetIndustries() ([]*domain.Industry, error) {
	industries, err := service.store.GetIndustries()
	if err != nil {
		return nil, nil
	}
	return industries, nil
}

func (service *AdditionalUserService) GetDegrees() ([]*domain.Degree, error) {
	degrees, err := service.store.GetDegrees()
	if err != nil {
		return nil, nil
	}
	return degrees, nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}
