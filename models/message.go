package models

import "time"

type Message struct {
	ID         int    // to identify messages in history
	Body       string `json:"body"`
	Attachment []byte `json:"attachmet"`
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
