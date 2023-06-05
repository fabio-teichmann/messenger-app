package main

import (
	"fmt"
	"messenger-app/models"
	"sync"
	"time"
)

type Event struct {
	SubjectID int // on which queue to publish the message
	TargetID  int // which user should receive the message
	Data      models.Message
}

type Subscriber interface {
	NotifyCallback(Event)
}

type Subject interface {
	AddSubscriber(Subscriber)
	RemoveSubscriber(Subscriber)
	NotifySubscriber(Event)
}

type EventSubscriber struct {
	ID int
}

type EventSubject struct {
	ID        int
	Observers sync.Map
}

func (subscriber *EventSubscriber) NotifyCallback(event Event) {
	fmt.Printf("Sender: %v, Target: %v, Receiver: %v, Time: %v, Message %s\n", event.SubjectID, event.TargetID, subscriber.ID, event.Data.Time, event.Data.Body)
}

func (subject *EventSubject) AddSubscriber(sub EventSubscriber) {
	subject.Observers.Store(sub, struct{}{})
}

func (subject *EventSubject) RemoveSubscriber(sub EventSubscriber) {
	subject.Observers.Delete(sub)
}

func (es *EventSubject) NotifySubscriber(event Event) {
	es.Observers.Range(func(key interface{}, value interface{}) bool {
		fmt.Println(key)
		if key == nil {
			fmt.Printf("could not find matching Subscriber with id %v to event id: %v", event.TargetID, event.SubjectID)
			return false
		}
		es := key.(EventSubscriber)

		if es.ID == event.TargetID {
			// found matching subscriber
			es.NotifyCallback(event)
			return false
		}
		return true
	})
	// fmt.Printf("could not find matching Subscriber with id %v to event id: %v\n", event.TargetID, event.SubjectID)
}

// func SendMessage(sender EventSubscriber, target EventSubject, message models.Message) {
// 	event := Event{
// 		SubjectID: sender.ID,
// 		TargetID:  target.ID,
// 		Data:      message,
// 	}

// }

func main() {
	var user1 = EventSubscriber{ID: 1}
	var user2 = EventSubscriber{ID: 2}
	var user3 = EventSubscriber{ID: 3}

	var chat12 = EventSubject{
		ID:        12,
		Observers: sync.Map{},
	}
	var chat13 = EventSubject{
		ID:        13,
		Observers: sync.Map{},
	}

	// subscribing users to chats
	chat12.AddSubscriber(user2)
	chat12.AddSubscriber(user1)

	chat13.AddSubscriber(user3)
	chat13.AddSubscriber(user1)

	for _, chat := range []*EventSubject{&chat13, &chat12, &chat13} {
		message := models.Message{
			Body: fmt.Sprintf("%s_%v", "test message", chat.ID),
			Time: time.Now(),
		}

		event := Event{
			SubjectID: user1.ID,
			TargetID:  chat.ID % 10,
			Data:      message,
		}
		// fmt.Println(event)

		chat.NotifySubscriber(event)
		// user1.NotifyCallback(event)

	}
}
