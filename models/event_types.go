package models

// EventType
type EventType int

// event types
const (
	MSG_SENT = iota
	MSG_RECEIVED
	MSG_SEEN
	NEW_USER
	USER_LOGIN
	USER_ONLINE
	USER_IN_CHAT
	USER_TYPING
	USER_LOGOUT
	CREATE_CHAT
	DELETE_CHAT
)

type ControlMsg int

const (
	DoExit = iota
	ExitOK
)
