package models

// EventType
type EventType int

// event types
const (
	MSG_SENT = iota
	MSG_RECEIVED
	MSG_SEEN
	NEW_USER
	USER_ONLINE
	USER_IN_CHAT
	USER_TYPING
)
