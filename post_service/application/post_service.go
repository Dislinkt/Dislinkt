package application

import (
	"context"
	"errors"
	"github.com/dislinkt/common/tracer"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/domain"
)

type PostService struct {
	store                      domain.PostStore
	createJobOfferOrchestrator *CreateJobOfferOrchestrator
}

func NewPostService(store domain.PostStore, createJobOfferOrchestrator *CreateJobOfferOrchestrator) *PostService {
	return &PostService{
		store:                      store,
		createJobOfferOrchestrator: createJobOfferOrchestrator,
	}
}

func (service *PostService) Get(ctx context.Context, id primitive.ObjectID) (*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "Get-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.Get(id)
}

func (service *PostService) GetAll(ctx context.Context) ([]*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAll()
}

func (service *PostService) Insert(ctx context.Context, post *domain.Post) error {
	span := tracer.StartSpanFromContext(ctx, "Insert-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.Insert(post)
}

func (service *PostService) GetAllByUserId(ctx context.Context, uuid string) ([]*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAllByUserId(uuid)
}

func (service *PostService) GetAllByConnectionIds(ctx context.Context, uuids []string) ([]*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByConnectionIds-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAllByConnectionIds(uuids)
}

func (service *PostService) CreateComment(ctx context.Context, post *domain.Post, comment *domain.Comment) error {
	span := tracer.StartSpanFromContext(ctx, "CreateComment-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.CreateComment(post, comment)
}

func (service *PostService) LikePost(ctx context.Context, post *domain.Post, userId string) error {
	span := tracer.StartSpanFromContext(ctx, "LikePost-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.LikePost(post, userId)
}

func (service *PostService) DislikePost(ctx context.Context, post *domain.Post, userId string) error {
	span := tracer.StartSpanFromContext(ctx, "DislikePost-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.DislikePost(post, userId)
}

func (service *PostService) GetRecent(ctx context.Context, uuid string) ([]*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRecent-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetRecent(uuid)
}

func (service *PostService) GetAllJobOffers(ctx context.Context) ([]*domain.JobOffer, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllJobOffers-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAllJobOffers()
}

func (service *PostService) InsertJobOfferOrc(ctx context.Context, offer *domain.JobOffer) error {
	span := tracer.StartSpanFromContext(ctx, "InsertJobOfferOrch-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err := service.createJobOfferOrchestrator.Start(offer)
	if err != nil {
		return err
	}
	return err
}

func (service *PostService) InsertJobOffer(ctx context.Context, offer *domain.JobOffer) error {
	span := tracer.StartSpanFromContext(ctx, "InsertJobOffer-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err := service.store.InsertJobOffer(offer)
	if err != nil {
		return err
	}
	return err
}

func (service *PostService) DeleteJobOffer(ctx context.Context, jobOffer *domain.JobOffer) error {
	span := tracer.StartSpanFromContext(ctx, "DeleteJobOffer-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.DeleteJobOffer(jobOffer)
}

func (service *PostService) InsertUser(ctx context.Context, user *domain.User) error {
	span := tracer.StartSpanFromContext(ctx, "InsertUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return errors.New("Invalid user data")
	}
	return service.store.InsertUser(user)
}

func (service *PostService) DeleteUser(ctx context.Context, user *domain.User) error {
	span := tracer.StartSpanFromContext(ctx, "DeleteUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return errors.New("Invalid user data")
	}
	return service.store.DeleteUser(user)
}

func (service *PostService) UpdateUser(ctx context.Context, user *domain.User) error {
	span := tracer.StartSpanFromContext(ctx, "UpdateUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if err := validator.New().Struct(user); err != nil {
		//	logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return errors.New("Invalid user data")
	}
	return service.store.UpdateUser(user)
}

func (service *PostService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetUser-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetUser(id)
}

func (service *PostService) SearchJobOffers(ctx context.Context, searchText string) ([]*domain.JobOffer, error) {
	span := tracer.StartSpanFromContext(ctx, "SearchJobOffers-Service")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.SearchJobOffers(searchText)
}
