package models

import "errors"

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProfilePic bool   `json:"profile_pic"`
}

func (user *User) CreateEventMessage(message *Message, target *User) (*Event, error) {
	if message == nil {
		return nil, errors.New("missing message")
	}
	if target == nil {
		return nil, errors.New("no target defined")
	}

	event := &Event{
		SubjectID: user.ID,
		TargetID:  target.ID,
		Data:      *message,
	}
	return event, nil
}
