package models

import (
	"errors"
	"fmt"
)

type Subscriber interface {
	NotifyCallback(*Event)
	CreateEvent(EventType, Message, *EventSubscriber) (*Event, error)
}

type User struct {
	ID         int    `json:"id" bson:"id"`
	Name       string `json:"name" bson:"name"`
	ProfilePic bool   `json:"profile_pic" bson:"profile_pic,omitempty"`
	// Chats      []Chat // list of conversations
}

type EventSubscriber struct {
	User User
}

func NewEventSubscriber(user User) *EventSubscriber {
	return &EventSubscriber{User: user}
}

func (subscriber *EventSubscriber) NotifyCallback(event *Event) {
	// fmt.Printf("EventType: %v, \n", event.EventType)
	fmt.Printf("Sender: %v, Target: %v, Receiver: %v, Time: %v, Message %s\n", event.Sender.User.ID, event.Target.User.ID, subscriber.User.ID, event.Data.Time, event.Data.Body)

	if event.EventType == MSG_SENT {
		fmt.Printf("Event: MSG_SENT, Sender: %v, Target: %v\n", event.Sender, event.Target)
		fmt.Printf("initiate MSG_RECEIVED event...\n")
		// initiate MsgReceived
		e, err := subscriber.CreateEvent(MSG_RECEIVED, &event.Data, &event.Sender)
		if err != nil {
			fmt.Println(err)
		}
		// Need ESB to accept event
		fmt.Println(*e)

	} else if event.EventType == MSG_RECEIVED {
		fmt.Printf("Event: MSG_RECEIVED, Sender: %v, Target: %v\n", event.Sender, event.Target)
		// update message as Received
		event.Data.Received = true
	}
}

func (es *EventSubscriber) CreateEvent(eventType EventType, message *Message, target *EventSubscriber) (*Event, error) {
	if eventType == MSG_SENT && message == nil {
		return nil, errors.New("missing message")
	}
	if target == nil {
		return nil, errors.New("no target defined")
	}

	return &Event{
		SubjectID: 0,
		Sender:    *es,
		Target:    *target,
		Data:      *message,
		EventType: eventType,
	}, nil
}

// func (user *User) CreateEventMessage(message *Message, target *User) (*Event, error) {
// 	if message == nil {
// 		return nil, errors.New("missing message")
// 	}
// 	if target == nil {
// 		return nil, errors.New("no target defined")
// 	}

// 	event := &Event{
// 		Sender: *user,
// 		Target: *target,
// 		Data:   *message,
// 	}
// 	return event, nil
// }

// func (user *User) CreateEvent(eventType EventType, message *Message, target *User) (*Event, error) {
// 	if eventType == MSG_SENT && message == nil {
// 		return nil, errors.New("missing message")
// 	}
// 	if target == nil {
// 		return nil, errors.New("no target defined")
// 	}

// 	return &Event{
// 		SubjectID: 0,
// 		Sender:    *user,
// 		Target:    *target,
// 		Data:      *message,
// 		EventType: eventType,
// 	}, nil
// }
