package models

import "errors"

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProfilePic bool   `json:"profile_pic"`
	// Chats      []Chat // list of conversations
}

func (user *User) CreateEventMessage(message *Message, target *User) (*Event, error) {
	if message == nil {
		return nil, errors.New("missing message")
	}
	if target == nil {
		return nil, errors.New("no target defined")
	}

	event := &Event{
		Sender: *user,
		Target: *target,
		Data:   *message,
	}
	return event, nil
}

func (user *User) CreateEvent(eventType EventType, message *Message, target *User) (*Event, error) {
	if eventType == MSG_SENT && message == nil {
		return nil, errors.New("missing message")
	}
	if target == nil {
		return nil, errors.New("no target defined")
	}

	return &Event{
		SubjectID: 0,
		Sender:    *user,
		Target:    *target,
		Data:      *message,
		EventType: eventType,
	}, nil
}
