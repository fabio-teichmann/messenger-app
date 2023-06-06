package models

type Event struct {
	SubjectID int // on which queue to publish the message
	TargetID  int // which user should receive the message
	Data      Message
}

func (e *Event) SendToChat(chat Chat) {
	chat.Chat <- *e
	chat.History = append(chat.History, e.Data)
}
