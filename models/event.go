package models

// EventType
type EventType int

// event types
const (
	UserOnline = iota
	UserTyping
	MsgSent
	MsgReceived
)

type Event struct {
	SubjectID int       // on which queue to publish the message
	TargetID  int       // which user should receive the message
	Data      Message   // contains payload
	EventType EventType // to classify events
}

func (e *Event) SendToChat(chat Chat) {
	chat.Chat <- *e
	chat.History = append(chat.History, e.Data)
}
