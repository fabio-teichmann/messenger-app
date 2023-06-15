package models

import (
	"messenger-app/util"
	"time"
)

type Message struct {
	ID         uint32 `bson:"message_id"` // to identify messages in history
	Body       string `json:"body"`
	Attachment []byte `json:"attachmet" bson:"omitempty"`
	Time       time.Time
	Sent       bool
	Received   bool
}

func NewMessage(body string) Message {
	hash := util.CreateHash([]byte(body))

	message := Message{
		ID:   hash,
		Body: body,
		Time: time.Now(),
	}
	return message
}
