package events

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type JobOffer struct {
	Id            primitive.ObjectID
	Position      string
	Description   string
	Preconditions string
	DatePosted    time.Time
	Duration      int
	Location      string
	Title         string
	Field         string
}

type CreateJobOfferCommandType int8

const (
	CreateJobOfferInPost CreateJobOfferCommandType = iota
	RollbackJobOfferInPost
	CreateJobOfferInGraph
	CancelJobOfferCreation
	ApproveCreation
	UnknownCreateJobOfferCommand
)

type CreateJobOfferCommand struct {
	JobOffer JobOffer
	Type     CreateJobOfferCommandType
}

type CreateJobOfferReplyType int8

const (
	PostServiceCreated CreateJobOfferReplyType = iota
	PostServiceNotCreated
	PostServiceRolledBack
	ConnectionServiceCreated
	ConnectionServiceNotCreated
	UnknownCreateJobOfferReply
)

type CreateJobReply struct {
	JobOffer JobOffer
	Type     CreateJobOfferReplyType
}
