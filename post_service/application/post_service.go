package application

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/domain"
)

type PostService struct {
	store domain.PostStore
}

func NewPostService(store domain.PostStore) *PostService {
	return &PostService{store: store}
}

func (service *PostService) Get(id primitive.ObjectID) (*domain.Post, error) {
	return service.store.Get(id)
}

func (service *PostService) GetAll() ([]*domain.Post, error) {
	return service.store.GetAll()
}

func (service *PostService) Insert(post *domain.Post) error {
	return service.store.Insert(post)
}

func (service *PostService) GetAllByUserId(uuid string) ([]*domain.Post, error) {
	return service.store.GetAllByUserId(uuid)
}

func (service *PostService) GetAllByConnectionIds(uuids []string) ([]*domain.Post, error) {
	return service.store.GetAllByConnectionIds(uuids)
}

func (service *PostService) CreateComment(post *domain.Post, comment *domain.Comment) error {
	return service.store.CreateComment(post, comment)
}
