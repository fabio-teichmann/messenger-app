package models

import (
	"fmt"
	"sync"
)

type EventSubject struct {
	ID int
	// Queue     chan Event
	Observers sync.Map
}

// type EventSubscriber struct {
// 	User User
// }

type Chat struct {
	// Subscribers []*EventSubscriber
	Subject *EventSubject
	History []Message
	Chat    chan Event
}

var idCount = 0

// func (subscriber *EventSubscriber) NotifyCallback(event *Event) {
// 	// fmt.Printf("EventType: %v, \n", event.EventType)
// 	fmt.Printf("Sender: %v, Target: %v, Receiver: %v, Time: %v, Message %s\n", event.Sender.User.ID, event.Target.User.ID, subscriber.User.ID, event.Data.Time, event.Data.Body)

// 	if event.EventType == MSG_SENT {
// 		fmt.Printf("Event: MSG_SENT, Sender: %v, Target: %v\n", event.Sender, event.Target)
// 		fmt.Printf("initiate MSG_RECEIVED event...\n")
// 		// initiate MsgReceived
// 		e, err := subscriber.CreateEvent(MSG_RECEIVED, &event.Data, &event.Sender)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		// Need ESB to accept event
// 		fmt.Println(*e)

// 	} else if event.EventType == MSG_RECEIVED {
// 		fmt.Printf("Event: MSG_RECEIVED, Sender: %v, Target: %v\n", event.Sender, event.Target)
// 		// update message as Received
// 		event.Data.Received = true
// 	}
// }

// Creates a chat queue
func (subscriber *EventSubscriber) CreateChat(subscribers []*EventSubscriber) Chat {

	// create a subject
	subject := &EventSubject{
		ID:        idCount,
		Observers: sync.Map{},
	}
	idCount++

	// subscribe all subscribers to subject
	for _, sub := range append(subscribers, subscriber) {
		// fmt.Println("CreateChat:", sub)
		subject.AddSubscriber(sub)
	}

	return Chat{
		// Subscribers: subs,
		Subject: subject,
		History: []Message{},
		Chat:    make(chan Event),
	}
}

type ControlMsg int

const (
	DoExit = iota
	ExitOK
)

func (chat *Chat) ReadMessages(control chan ControlMsg) {

	for {
		select {
		case msg := <-control:
			switch msg {
			case DoExit:
				fmt.Println("exit read message")
				control <- ExitOK
				return
			}
		case message := <-chat.Chat:
			chat.Subject.NotifySubscriber(&message)
		}
	}
}
