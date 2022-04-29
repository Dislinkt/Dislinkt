package application

import (
	"fmt"

	"github.com/dislinkt/additional_user_service/domain"
)

type AdditionalUserService struct {
	store domain.AdditionalUserStore
}

func NewAdditionalUserService(store domain.AdditionalUserStore) *AdditionalUserService {
	return &AdditionalUserService{
		store: store,
	}
}

func (service *AdditionalUserService) CreateEducation(uuid string, education *domain.Education) (*domain.Education,
	error) {

	_, err := service.store.FindOrCreateDocument(uuid)
	if err != nil {
		return nil, err
	}

	insertEducation, err := service.store.InsertEducation(uuid, education)
	if err != nil {
		return nil, err
	}

	return insertEducation, nil
}

func (service *AdditionalUserService) FindUserEducation(userUUID string) (user *[]domain.Education, err error) {
	document, err := service.store.FindUserDocument(userUUID)
	fmt.Println(document.Educations)
	if err != nil {
		return nil, err
	}
	return &document.Educations, nil
}
