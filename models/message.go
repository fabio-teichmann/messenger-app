package models

import "time"

type Message struct {
	ID         int    `bson:"message_id"` // to identify messages in history
	Body       string `json:"body"`
	Attachment []byte `json:"attachmet" bson:"omitempty"`
	Time       time.Time
	Sent       bool
	Received   bool
}

func NewMessage(body string) Message {
	message := Message{
		Body: body,
		Time: time.Now(),
	}
	return message
}
