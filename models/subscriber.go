package models

import (
	"context"
	"errors"
	"fmt"
	"messenger-app/util"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscriber interface {
	NotifyCallback(*Event)
	CreateEvent(EventType, Message, *EventSubscriber) (*Event, error)
}

type User struct {
	ID         uint32 `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	ProfilePic bool   `json:"profile_pic" bson:"profile_pic,omitempty"`
	// Chats      []Chat // list of conversations
}

type EventSubscriber struct {
	User
	Chats []EventSubscriber
}

func NewEventSubscriber(user User) *EventSubscriber {
	return &EventSubscriber{User: user}
}

func NewEventSubscriberByName(userName string) *EventSubscriber {
	hash := util.CreateHash([]byte(userName))
	return &EventSubscriber{
		User{
			ID:   hash,
			Name: userName,
		},
		[]EventSubscriber{},
	}
}

func NewEventSubscriberWithEvent(userName string) (*EventSubscriber, *Event) {
	es := NewEventSubscriberByName(userName)

	event, err := es.CreateEvent(NEW_USER, &Message{Body: "new user"}, nil)
	if err != nil {
		fmt.Println("error occurred creating event:", err)
		return es, nil
	}

	return es, event
}

func (subscriber *EventSubscriber) NotifyCallback(ctx context.Context, ac *AppControler, event *Event) {
	// fmt.Printf("EventType: %v, \n", event.EventType)
	fmt.Printf("Sender: %v, Target: %v, Receiver: %v, Message %s\n", event.Sender.User.ID, event.Target.User.ID, subscriber.User.ID, event.Data.Body)

	if event.SubjectID == MSG_SENT {
		fmt.Printf("Event: MSG_SENT, Sender: %v, Target: %v\n", event.Sender, event.Target)
		// set message to SENT
		err := ac.SetMessageToSent(ctx, event.Data.ID)
		if err != nil {
			fmt.Println(err)
		}

		// initiate MsgReceived
		fmt.Printf("initiate MSG_RECEIVED event...\n")
		e, err := event.Target.CreateEvent(MSG_RECEIVED, &event.Data, &event.Sender)
		if err != nil {
			fmt.Println(err)
		}

		ac.AcceptEvent(ctx, e)

	} else if event.SubjectID == MSG_RECEIVED {
		fmt.Printf("Event: MSG_RECEIVED, Sender: %v, Target: %v\n", event.Sender, event.Target)
		// update message as Received
		err := ac.SetMessageToRcvd(ctx, event.Data.ID)
		if err != nil {
			fmt.Println(err)
		}

	} else if event.SubjectID == USER_LOGIN {
		fmt.Printf("User Login: %v", event.Sender)

		// trigger USER_ONLINE --> need to notify all subscribers that chat with event.Sender
		event.Sender.CreateEvent(USER_ONLINE, nil, nil)
	}
}

func (es *EventSubscriber) CreateEvent(eventType EventType, message *Message, target *EventSubscriber) (*Event, error) {
	if message == nil {
		return nil, errors.New("missing message")
	}

	event := Event{
		ID:        primitive.NewObjectID(),
		SubjectID: eventType,
		Sender:    *es,
		Data:      *message,
		Time:      time.Now(),
	}

	if eventType == NEW_USER {
		return &event, nil
	}

	if target == nil {
		return nil, errors.New("no target defined")
	}
	event.Target = *target
	return &event, nil
}
