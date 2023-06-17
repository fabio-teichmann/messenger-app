package models

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	User
}

func NewEventSubscriber(user User) *EventSubscriber {
	return &EventSubscriber{User: user}
}

func (subscriber *EventSubscriber) NotifyCallback(ctx context.Context, ac *AppControler, event *Event) {
	// fmt.Printf("EventType: %v, \n", event.EventType)
	fmt.Printf("Sender: %v, Target: %v, Receiver: %v, Message %s\n", event.Sender.User.ID, event.Target.User.ID, subscriber.User.ID, event.Data.Body)

	if event.SubjectID == MSG_SENT {
		fmt.Printf("Event: MSG_SENT, Sender: %v, Target: %v\n", event.Sender, event.Target)
		fmt.Printf("initiate MSG_RECEIVED event...\n")
		// initiate MsgReceived
		fmt.Println(event)
		err := ac.UpdateEventMessageToSentById(ctx, event.ID)
		if err != nil {
			fmt.Println(err)
		}
		event.SubjectID = MSG_RECEIVED
		go func() { ac.MsgRcvd.Queue <- *event }()

	} else if event.SubjectID == MSG_RECEIVED {
		fmt.Printf("Event: MSG_RECEIVED, Sender: %v, Target: %v\n", event.Sender, event.Target)
		// update message as Received
		err := ac.UpdateEventMessageToRcvdById(ctx, event.ID)
		if err != nil {
			fmt.Println(err)
		}
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
		ID:        primitive.NewObjectID(),
		SubjectID: eventType,
		Sender:    *es,
		Target:    *target,
		Data:      *message,
		// EventType: eventType,
	}, nil
}
