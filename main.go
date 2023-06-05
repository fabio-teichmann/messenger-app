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

type Chat struct {
	// Subscribers []*EventSubscriber
	Subject *EventSubject
	History []models.Message
	Chat    chan Event
}

var idCount = 1

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
		fmt.Println(es.ID, event.Data, key)
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
		subject.AddSubscriber(*sub)
	}

	// var subs []*EventSubscriber
	// subs = append(subs, subscriber)

	return Chat{
		// Subscribers: subs,
		Subject: subject,
		History: []models.Message{},
		Chat:    make(chan Event),
	}
}

func (e *Event) SendToChat(chat Chat) {
	go func() {
		// defer close(chat.Chat)
		chat.Chat <- *e
		chat.History = append(chat.History, e.Data)
		// close(chat.Chat)
	}()
}

func (chat *Chat) ReadMessages() {
	for message := range chat.Chat {
		chat.Subject.NotifySubscriber(message)
	}
}

func main() {
	var user1 = EventSubscriber{ID: 1}
	var user2 = EventSubscriber{ID: 2}
	// var user3 = EventSubscriber{ID: 3}

	var chat12 = user1.CreateChat([]*EventSubscriber{&user2})
	// var chat13 = user1.CreateChat([]*EventSubscriber{&user3})

	for _, i := range []int{1, 2, 3, 1, 2, 3} {
		message := models.Message{
			Body: fmt.Sprintf("%s_%v", "test message", i),
			Time: time.Now(),
		}

		event := Event{
			SubjectID: chat12.Subject.ID,
			TargetID:  i,
			Data:      message,
		}

		event.SendToChat(chat12)
	}
	chat12.ReadMessages()
	// chat13.ReadMessages()
	// var chat12 = EventSubject{
	// 	ID:        12,
	// 	Observers: sync.Map{},
	// }
	// var chat13 = EventSubject{
	// 	ID:        13,
	// 	Observers: sync.Map{},
	// }

	// subscribing users to chats
	// chat12.AddSubscriber(user2)
	// chat12.AddSubscriber(user1)

	// chat13.AddSubscriber(user3)
	// chat13.AddSubscriber(user1)

	// for _, chat := range []*EventSubject{&chat13, &chat12, &chat13} {
	// 	message := models.Message{
	// 		Body: fmt.Sprintf("%s_%v", "test message", chat.ID),
	// 		Time: time.Now(),
	// 	}

	// 	event := Event{
	// 		SubjectID: user1.ID,
	// 		TargetID:  chat.ID % 10,
	// 		Data:      message,
	// 	}
	// 	// fmt.Println(event)

	// 	chat.NotifySubscriber(event)
	// user1.NotifyCallback(event)

	// }
}
