package models

// EventType
type EventType int

// event types
const (
	USER_ONLINE = iota
	USER_TYPING
	MSG_SENT
	MSG_RECEIVED
)

type Event struct {
	SubjectID int             // on which queue to publish the message
	Sender    EventSubscriber // event origin
	Target    EventSubscriber // event destination
	Data      Message         // contains payload
	EventType EventType       // to classify events
}

func (e *Event) SendToChat(chat Chat) {
	chat.Chat <- *e
	chat.History = append(chat.History, e.Data)
}
