package api

import (
	"context"
	"github.com/dislinkt/common/interceptor"
	logger "github.com/dislinkt/common/logging"
	"time"

	"post_service/domain"

	pb "github.com/dislinkt/common/proto/post_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/application"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service *application.PostService
	logger  *logger.Logger
}

func NewPostHandler(service *application.PostService) *PostHandler {
	logger := logger.InitLogger(context.TODO())
	return &PostHandler{
		service: service,
		logger:  logger,
	}
}

func (handler PostHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.logger.WarnLogger.Warn("PNF {%s}", ctx.Value(interceptor.LoggedInUserKey{}).(string))
		return nil, err
	}
	postPb := mapPost(post)
	response := &pb.GetResponse{Post: postPb}
	return response, nil
}

func (handler *PostHandler) GetRecent(ctx context.Context, request *pb.GetRequest) (*pb.GetMultipleResponse, error) {
	id := request.Id
	posts, err := handler.service.GetRecent(id)
	if err != nil {
		return nil, err
	}
	response := &pb.GetMultipleResponse{Posts: []*pb.Post{}}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	return response, nil
}

func (handler *PostHandler) GetAllByUserId(ctx context.Context, request *pb.GetRequest) (*pb.GetMultipleResponse, error) {
	id := request.Id
	posts, err := handler.service.GetAllByUserId(id)
	if err != nil {
		return nil, err
	}
	response := &pb.GetMultipleResponse{Posts: []*pb.Post{}}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	return response, nil
}

/*
func (handler *PostHandler) GetAllByConnectionIds(ctx context.Context, request *pb.GetRequest) (*pb.GetMultipleResponse, error) {
	id := request.Id
	posts, err := handler.service.GetAllByConnectionIds(id)
	if err != nil {
		return nil, err
	}
	response := &pb.GetMultipleResponse{Posts: []*pb.Post{}}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	return response, nil
}
*/

func (handler *PostHandler) GetAll(ctx context.Context, request *pb.Empty) (*pb.GetMultipleResponse, error) {
	posts, err := handler.service.GetAll()
	if err != nil {
		return nil, err
	}
	response := &pb.GetMultipleResponse{Posts: []*pb.Post{}}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	return response, nil
}

func (handler *PostHandler) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (*pb.Empty, error) {
	handler.logger.InfoLogger.Infof("POST rr: PC {%s}", ctx.Value(interceptor.LoggedInUserKey{}).(string))
	post := mapNewPost(request.Post)
	err := handler.service.Insert(post)
	if err != nil {
		handler.logger.WarnLogger.Warn("WPC {%s}", ctx.Value(interceptor.LoggedInUserKey{}).(string))
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (handler *PostHandler) CreateComment(ctx context.Context, request *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(request.PostId)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	comment := mapNewComment(request.Comment)
	err = handler.service.CreateComment(post, comment)
	if err != nil {
		handler.logger.WarnLogger.Warn("WCC {%s}", ctx.Value(interceptor.LoggedInUserKey{}).(string))
		return nil, err
	}

	return &pb.CreateCommentResponse{
		Comment: request.Comment,
	}, nil
}

func (handler *PostHandler) LikePost(ctx context.Context, request *pb.ReactionRequest) (*pb.Empty, error) {
	objectId, err := primitive.ObjectIDFromHex(request.PostId)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	err = handler.service.LikePost(post, request.UserId)
	if err != nil {
		handler.logger.WarnLogger.Warn("WR {%s}", ctx.Value(interceptor.LoggedInUserKey{}).(string))
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (handler *PostHandler) DislikePost(ctx context.Context, request *pb.ReactionRequest) (*pb.Empty, error) {
	objectId, err := primitive.ObjectIDFromHex(request.PostId)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	err = handler.service.DislikePost(post, request.UserId)
	if err != nil {
		handler.logger.WarnLogger.Warn("WR {%s}", ctx.Value(interceptor.LoggedInUserKey{}).(string))
		return nil, err
	}

	return &pb.Empty{}, nil
}

/* JOB OFFERS */

func (handler *PostHandler) GetAllJobOffers(ctx context.Context, request *pb.SearchMessage) (*pb.GetAllJobOffers, error) {
	var offers []*domain.JobOffer
	var err error

	if len(request.SearchText) == 0 {
		offers, err = handler.service.GetAllJobOffers()
	} else {
		offers, err = handler.service.SearchJobOffers(request.SearchText)
	}

	if err != nil {
		return nil, err
	}
	response := &pb.GetAllJobOffers{JobOffers: []*pb.JobOffer{}}
	for _, offer := range offers {
		y, m, d := offer.DatePosted.Date()
		endDate := time.Date(y, m, d+offer.Duration, 0, 0, 0, 0, time.Local)

		if time.Now().Before(endDate) {
			current := mapJobOffer(offer)
			response.JobOffers = append(response.JobOffers, current)
		}
	}

	return response, nil
}

func (handler *PostHandler) CreateJobOffer(ctx context.Context, request *pb.CreateJobOfferRequest) (*pb.Empty, error) {
	handler.logger.InfoLogger.Infof("POST rr: JOC {%s}", ctx.Value(interceptor.LoggedInUserKey{}).(string))
	offer := mapNewJobOffer(request.JobOffer)
	err := handler.service.InsertJobOffer(offer)
	if err != nil {
		handler.logger.WarnLogger.Warn("WJOC {%s}", ctx.Value(interceptor.LoggedInUserKey{}).(string))
		return nil, err
	}
	return &pb.Empty{}, nil
}

/* REACTIONS I COMMENTS */

func (handler *PostHandler) GetAllLikesForPost(ctx context.Context, request *pb.GetRequest) (*pb.GetReactionsResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}

	response := &pb.GetReactionsResponse{Users: []*pb.User{}}
	for _, reaction := range post.Reactions {
		if reaction.Reaction == domain.LIKED {
			user, err := handler.service.GetUser(reaction.UserId)
			if err != nil {
				return nil, err
			}
			current := mapUserReaction(user)
			response.Users = append(response.Users, current)
		}
	}

	if err != nil {
		return nil, err
	}
	return response, nil
}

func (handler *PostHandler) GetAllDislikesForPost(ctx context.Context, request *pb.GetRequest) (*pb.GetReactionsResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}

	response := &pb.GetReactionsResponse{Users: []*pb.User{}}
	for _, reaction := range post.Reactions {
		if reaction.Reaction == domain.DISLIKED {
			user, err := handler.service.GetUser(reaction.UserId)
			if err != nil {
				return nil, err
			}
			current := mapUserReaction(user)
			response.Users = append(response.Users, current)
		}
	}

	if err != nil {
		return nil, err
	}
	return response, nil
}

func (handler *PostHandler) GetAllCommentsForPost(ctx context.Context, request *pb.GetRequest) (*pb.GetAllCommentsResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}

	response := &pb.GetAllCommentsResponse{Comments: []*pb.Comment{}}
	for _, comment := range post.Comments {
		user, err := handler.service.GetUser(comment.UserId)
		if err != nil {
			return nil, err
		}
		current := mapUserCommentsForPost(user, comment.CommentText)
		response.Comments = append(response.Comments, current)
	}

	if err != nil {
		return nil, err
	}
	return response, nil
}
