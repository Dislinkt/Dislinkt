package api

import (
	"context"
	"fmt"
	"github.com/dislinkt/common/interceptor"
	connectionGw "github.com/dislinkt/common/proto/connection_service"
	notificationGw "github.com/dislinkt/common/proto/notification_service"
	userGw "github.com/dislinkt/common/proto/user_service"
	"post_service/infrastructure/persistence"
	"time"

	"post_service/domain"

	pb "github.com/dislinkt/common/proto/post_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post_service/application"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service *application.PostService
}

func NewPostHandler(service *application.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (handler PostHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	postPb := mapPost(post)
	response := &pb.GetResponse{Post: postPb}
	return response, nil
}

func (handler *PostHandler) GetRecent(ctx context.Context, request *pb.GetRequest) (*pb.GetMultipleResponse, error) {
	id := request.Id

	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	privacyResponse, _ := persistence.UserClient("user_service:8000").CheckIfUserIsPrivate(context.TODO(), &userGw.GetOneMessage{Id: id})
	isPrivate := privacyResponse.IsPrivate
	areUsersConnected := false
	isUserTheSame := false
	if username != "" {
		userResponse, _ := persistence.UserClient("user_service:8000").GetUserByUsername(context.TODO(), &userGw.GetOneByUsernameMessage{Username: username})
		connectionResponse, _ := persistence.ConnectionClient("connection_service:8000").CheckIfUsersConnected(context.TODO(), &connectionGw.CheckConnection{Uuid1: userResponse.User.Id, Uuid2: id})
		areUsersConnected = connectionResponse.IsConnected
		if id == userResponse.User.Id {
			isUserTheSame = true
		}
	}
	var posts []*domain.Post
	var err error
	response := &pb.GetMultipleResponse{Posts: []*pb.Post{}}

	if areUsersConnected || (!areUsersConnected && !isPrivate) || isUserTheSame {
		posts, err = handler.service.GetRecent(id)
		if err != nil {
			return nil, err
		}
		response = &pb.GetMultipleResponse{Posts: []*pb.Post{}}
		for _, post := range posts {
			current := mapPost(post)
			response.Posts = append(response.Posts, current)
		}
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
	username := fmt.Sprintf(ctx.Value(interceptor.LoggedInUserKey{}).(string))
	userResponse, _ := persistence.UserClient("user_service:8000").GetUserByUsername(context.TODO(), &userGw.GetOneByUsernameMessage{Username: username})
	post := mapNewPost(request.Post)
	post.UserId = userResponse.User.Id
	err := handler.service.Insert(post)
	if err != nil {
		return nil, err
	}
	_, _ = persistence.NotificationClient("notification_service:8000").SaveNotification(context.TODO(),
		&notificationGw.SaveNotificationRequest{Notification: mapNotification(username), UserId: userResponse.User.Id})
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

/* REACTIONS I COMMENTS */

func (handler *PostHandler) CreateJobOffer(ctx context.Context, request *pb.CreateJobOfferRequest) (*pb.Empty, error) {
	offer := mapNewJobOffer(request.JobOffer)
	err := handler.service.InsertJobOffer(offer)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

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
