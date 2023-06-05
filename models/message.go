package models

import "time"

type Message struct {
	Body       string `json:"body"`
	Attachment []byte `json:"attachmet"`
	Time       time.Time
}
