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
	SubjectID int             `bson:"subject_id"`  // on which queue to publish the message
	Sender    EventSubscriber `bson:"sender"`      // event origin
	Target    EventSubscriber `bson:"target"`      // event destination
	Data      Message         `bson:"data,inline"` // contains payload
	EventType EventType       `bson:"event_type"`  // to classify events
}

func (e *Event) SendToChat(chat Chat) {
	chat.Chat <- *e
	chat.History = append(chat.History, e.Data)
}
