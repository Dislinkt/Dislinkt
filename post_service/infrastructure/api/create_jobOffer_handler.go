package api

import (
	"fmt"
	"github.com/dislinkt/common/saga/events"
	saga "github.com/dislinkt/common/saga/messaging"
	"post_service/application"
)

type CreateJobOfferCommandHandler struct {
	postService       *application.PostService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewJobOfferCommandHandler(postService *application.PostService, publisher saga.Publisher,
	subscriber saga.Subscriber) (*CreateJobOfferCommandHandler, error) {
	o := &CreateJobOfferCommandHandler{
		postService:       postService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateJobOfferCommandHandler) handle(command *events.CreateJobOfferCommand) {
	reply := events.CreateJobReply{JobOffer: command.JobOffer}

	switch command.Type {
	case events.CreateJobOfferInPost:
		fmt.Println("posthandler-createJobOffer")
		jobOffer := mapPostCommandCreateJob(command)

		err := handler.postService.InsertJobOffer(jobOffer)
		if err != nil {
			fmt.Println("posthandler-createJobOfferError")
			reply.Type = events.PostServiceNotCreated
			return
		}
		reply.Type = events.PostServiceCreated
	case events.RollbackJobOfferInPost:
		fmt.Println("posthandler-rollbackJobOffer")
		err := handler.postService.DeleteJobOffer(mapPostCommandCreateJob(command))
		if err != nil {
			return
		}
		reply.Type = events.PostServiceRolledBack
	default:
		reply.Type = events.UnknownCreateJobOfferReply
	}

	if reply.Type != events.UnknownCreateJobOfferReply {
		fmt.Println(reply)
		_ = handler.replyPublisher.Publish(reply)
	}
}
