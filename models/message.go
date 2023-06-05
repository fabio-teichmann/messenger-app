package models

type Message struct {
	Body       []byte `json:"body"`
	Attachment []byte `json:"attachmet"`
}
