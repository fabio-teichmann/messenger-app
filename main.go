package main

import (
	"fmt"
	"messenger-app/models"
	"sync"
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
	fmt.Printf("Sender: %v, Receiver: %v, Time: %v, Message %s\n", event.SubjectID, event.TargetID, event.Data.Time, event.Data.Body)
}

func (subject *EventSubject) AddSubscriber(sub EventSubscriber) {
	subject.Observers.Store(sub, struct{}{})
}

func (subject *EventSubject) RemoveSubscriber(sub EventSubscriber) {
	subject.Observers.Delete(sub)
}

func (es *EventSubject) NotifySubscriber(event Event) {
	es.Observers.Range(func(key interface{}, value interface{}) bool {

		if key == nil || value == nil {
			fmt.Printf("could not find matching Subscriber to event id: %v", event.SubjectID)
			return false
		}

		if key.(EventSubscriber).ID == event.TargetID {
			// found matching subscriber
			key.(EventSubscriber).NotifyCallback(event)
		}
	})
}

func main() {

}
