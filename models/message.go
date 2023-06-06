package models

import "time"

type Message struct {
	Body       string `json:"body"`
	Attachment []byte `json:"attachmet"`
	Time       time.Time
}

func NewMessage(body string) Message {
	message := Message{
		Body: body,
		Time: time.Now(),
	}
	return message
}
