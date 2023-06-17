package models

import (
	"messenger-app/util"
	"time"
)

type Message struct {
	ID         uint32    `bson:"message_id"` // to identify messages in history
	Body       string    `json:"body"`
	Attachment []byte    `json:"attachmet" bson:"omitempty"`
	Sent       bool      `bson:"sent"`
	TimeSent   time.Time `bson:"time_sent"`
	Received   bool      `bson:"received"`
	TimeRcvd   time.Time `bson:"time_rcvd"`
	Seen       bool      `bson:"seen"`
	TimeSeen   time.Time `bson:"time_seen"`
}

func NewMessage(body string) *Message {
	hash := util.CreateHash([]byte(body))

	message := Message{
		ID:   hash,
		Body: body,
	}
	return &message
}
