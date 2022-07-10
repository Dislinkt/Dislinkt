package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/dislinkt/additional_user_service/domain"
	"github.com/dislinkt/common/tracer"
	uuid "github.com/gofrs/uuid"
)

type AdditionalUserService struct {
	store                       domain.AdditionalUserStore
	addEducationOrchestrator    *AddEducationOrchestrator
	addSkillOrchestrator        *AddSkillOrchestrator
	deleteSkillOrchestrator     *DeleteSkillOrchestrator
	updateSkillOrchestrator     *UpdateSkillOrchestrator
	deleteEducationOrchestrator *DeleteEducationOrchestrator
	updateEducationOrchestrator *UpdateEducationOrchestrator
}

func (service *AdditionalUserService) CreateDocument(ctx context.Context, uuid string) error {
	span := tracer.StartSpanFromContext(ctx, "CreateDocument-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return errors.New("Invalid uuid")
	}
	_, err := service.store.CreateUserDocument(uuid)
	if err != nil {
		return err
	}

	return nil
}

func (service *AdditionalUserService) DeleteDocument(ctx context.Context, uuid string) error {
	span := tracer.StartSpanFromContext(ctx, "DeleteDocument-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func NewAdditionalUserService(store domain.AdditionalUserStore, addEducationOrchestrator *AddEducationOrchestrator,
	addSkillOrchestrator *AddSkillOrchestrator, deleteSkillOrchestrator *DeleteSkillOrchestrator,
	updateSkillOrchestrator *UpdateSkillOrchestrator, deleteEducationOrchestrator *DeleteEducationOrchestrator,
	updateEducationOrchestrator *UpdateEducationOrchestrator) *AdditionalUserService {
	return &AdditionalUserService{
		store:                       store,
		addEducationOrchestrator:    addEducationOrchestrator,
		addSkillOrchestrator:        addSkillOrchestrator,
		deleteSkillOrchestrator:     deleteSkillOrchestrator,
		updateSkillOrchestrator:     updateSkillOrchestrator,
		deleteEducationOrchestrator: deleteEducationOrchestrator,
		updateEducationOrchestrator: updateEducationOrchestrator,
	}
}

func (service *AdditionalUserService) CreateEducation(ctx context.Context, uuid string, education *domain.Education) (*domain.Education,
	error) {
	span := tracer.StartSpanFromContext(ctx, "CreateEducation-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) CreateEducationStart(ctx context.Context, uuid string, education *domain.Education) error {
	span := tracer.StartSpanFromContext(ctx, "CreateEducationStart-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err := service.addEducationOrchestrator.Start(education, uuid)
	if err != nil {
		return err
	}
	return err
}

func (service *AdditionalUserService) FindUserEducations(ctx context.Context, uuid string) (*map[string]domain.Education,
	error) {
	span := tracer.StartSpanFromContext(ctx, "FindUSerEducations-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Educations, nil
}

func (service *AdditionalUserService) UpdateUserEducation(ctx context.Context, uuid string, educationId string,
	education *domain.Education) (*map[string]domain.Education, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserEducation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) UpdateUserEducationStart(ctx context.Context, uuid string, educationId string,
	education *domain.Education) (*map[string]domain.Education, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserEducationStart")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	err := service.updateEducationOrchestrator.Start(uuid, educationId, education)
	if err != nil {
		return nil, err
	}
	userEducation, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userEducation.Educations, nil
}

func (service *AdditionalUserService) DeleteUserEducation(ctx context.Context, uuid string, additionID string) (*map[string]domain.Education, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteUserEducation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) DeleteUserEducationStart(ctx context.Context, uuid string, additionID string) (*map[string]domain.Education, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteUserEducationStart")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	err := service.deleteEducationOrchestrator.Start(uuid, additionID)
	userEducation, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userEducation.Educations, nil
}

// POSITION

func (service *AdditionalUserService) CreatePosition(ctx context.Context, uuid string, position *domain.Position) (*domain.Position,
	error) {
	span := tracer.StartSpanFromContext(ctx, "CreatePosition-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) FindUserPositions(ctx context.Context, uuid string) (*map[string]domain.Position, error) {
	span := tracer.StartSpanFromContext(ctx, "FindUSerPositions-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Positions, nil
}

func (service *AdditionalUserService) UpdateUserPosition(ctx context.Context, uuid string, positionId string,
	position *domain.Position) (*map[string]domain.Position, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserPosition-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) DeleteUserPosition(ctx context.Context, uuid string, additionID string) (*map[string]domain.Position,
	error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteUserPosition-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) CreateSkill(ctx context.Context, uuid string, skill *domain.Skill) (*domain.Skill,
	error) {
	span := tracer.StartSpanFromContext(ctx, "CreateSkill-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) CreateSkillStart(ctx context.Context, uuid string, skill *domain.Skill) error {
	span := tracer.StartSpanFromContext(ctx, "CreateSkillStart")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err := service.addSkillOrchestrator.Start(skill, uuid)
	if err != nil {
		return err
	}
	return err
}

func (service *AdditionalUserService) FindUserSkills(ctx context.Context, uuid string) (*map[string]domain.Skill, error) {
	span := tracer.StartSpanFromContext(ctx, "FindUserSkills-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Skills, nil
}

func (service *AdditionalUserService) FindUserSkill(ctx context.Context, skillId string, uuid string) (*domain.Skill, error) {
	span := tracer.StartSpanFromContext(ctx, "FindUserSkill-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	skills, err := service.FindUserSkills(ctx, uuid)
	mapSkill := *skills
	skill := mapSkill[skillId]
	if err != nil {
		return nil, err
	}
	fmt.Println("[AdditionalUserService: FindUserSkill]")
	fmt.Println("skill")
	fmt.Println(skill)
	return &skill, nil
}

func (service *AdditionalUserService) FindUserField(ctx context.Context, educationId string, uuid string) (*domain.Education, error) {
	span := tracer.StartSpanFromContext(ctx, "FindUserField-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	educations, err := service.FindUserEducations(ctx, uuid)
	mapEducations := *educations
	education := mapEducations[educationId]
	if err != nil {
		return nil, err
	}
	fmt.Println("[AdditionalUserService: FindUserField]")
	fmt.Println(education)
	return &education, nil
}

func (service *AdditionalUserService) UpdateUserSkill(ctx context.Context, uuid string, skillId string,
	skill *domain.Skill) (*map[string]domain.Skill, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserSkill-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) UpdateUserSkillStart(ctx context.Context, uuid string, skillId string,
	skill *domain.Skill) (*map[string]domain.Skill, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserSkillStart-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	err := service.updateSkillOrchestrator.Start(uuid, skillId, skill)
	if err != nil {
		return nil, err
	}
	userSkill, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userSkill.Skills, nil
}

func (service *AdditionalUserService) DeleteUserSkill(ctx context.Context, uuid string, additionID string) (*map[string]domain.Skill,
	error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteUSerSkill-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	err := service.store.DeleteUserSkill(additionID)
	userSkill, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}

	fmt.Println("[AdditionalUserService: DeleteUserSkill]")

	return userSkill.Skills, nil
}

func (service *AdditionalUserService) DeleteUserSkillStart(ctx context.Context, uuid string, additionID string) (*map[string]domain.Skill,
	error) {

	span := tracer.StartSpanFromContext(ctx, "DeleteUserSkillStart")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	err := service.deleteSkillOrchestrator.Start(uuid, additionID)
	if err != nil {
		return nil, err
	}
	userSkill, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return userSkill.Skills, nil
}

// INTEREST

func (service *AdditionalUserService) CreateInterest(ctx context.Context, uuid string, interest *domain.Interest) (*domain.Interest,
	error) {
	span := tracer.StartSpanFromContext(ctx, "CreateInterest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) FindUserInterests(ctx context.Context, uuid string) (*map[string]domain.Interest, error) {
	span := tracer.StartSpanFromContext(ctx, "FindUserInterests")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if !IsValidUUID(uuid) {
		return nil, errors.New("Invalid uuid")
	}

	document, err := service.store.FindUserDocument(uuid)
	if err != nil {
		return nil, err
	}
	return document.Interests, nil
}

func (service *AdditionalUserService) UpdateUserInterest(ctx context.Context, uuid string, interestId string,
	interest *domain.Interest) (*map[string]domain.Interest, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserInterest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) DeleteUserInterest(ctx context.Context, uuid string, additionID string) (*map[string]domain.Interest,
	error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteUserInterest-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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

func (service *AdditionalUserService) GetFieldOfStudies(ctx context.Context) ([]*domain.FieldOfStudy, error) {
	span := tracer.StartSpanFromContext(ctx, "GetFieldOfStudies-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	study, err := service.store.GetAllFieldOfStudy()
	if err != nil {
		return nil, nil
	}
	return study, nil
}

func (service *AdditionalUserService) GetSkills(ctx context.Context) ([]*domain.Skill, error) {
	span := tracer.StartSpanFromContext(ctx, "GetSkills-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	skills, err := service.store.GetSkills()
	if err != nil {
		return nil, nil
	}
	return skills, nil
}

func (service *AdditionalUserService) GetIndustries(ctx context.Context) ([]*domain.Industry, error) {
	span := tracer.StartSpanFromContext(ctx, "GetIndustries-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	industries, err := service.store.GetIndustries()
	if err != nil {
		return nil, nil
	}
	return industries, nil
}

func (service *AdditionalUserService) GetDegrees(ctx context.Context) ([]*domain.Degree, error) {
	span := tracer.StartSpanFromContext(ctx, "GetDegrees-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
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
