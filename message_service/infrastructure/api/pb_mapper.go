package api

import (
	pb "github.com/dislinkt/common/proto/message_service"
	"github.com/dislinkt/message_service/domain"
	"time"
)

func mapMessageHistory(messageHistory *domain.MessageHistory) *pb.MessageHistory {
	id := messageHistory.Id.Hex()

	messageHistoryPb := &pb.MessageHistory{
		Id:                   id,
		User1Id:              messageHistory.UserOneId,
		User2Id:              messageHistory.UserTwoId,
		UnreadMessagesNumber: 0,
	}
	for _, message := range messageHistory.Messages {
		messageHistoryPb.Messages = append(messageHistoryPb.Messages, &pb.Message{
			SenderId:    message.SenderId,
			ReceiverId:  message.ReceiverId,
			MessageText: message.MessageText,
			DateSent:    message.DateSent.String(),
			IsRead:      message.IsRead,
		})

		if !message.IsRead {
			messageHistoryPb.UnreadMessagesNumber++
		}
	}

	return messageHistoryPb
}

func mapNewMessage(messagePb *pb.Message) *domain.Message {
	message := &domain.Message{
		SenderId:    messagePb.SenderId,
		ReceiverId:  messagePb.ReceiverId,
		MessageText: messagePb.MessageText,
		DateSent:    time.Now(),
	}
	return message
}
